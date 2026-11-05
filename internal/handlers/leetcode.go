package handlers

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/rand/v2"
	"net/http"

	"github.com/bwmarrin/discordgo"
	"github.com/sudo-JP/tad/internal/db"
)

// handleOwedLeetcode expects args without the leading "owed", e.g. ["<@123>"].
func handleOwedLeetcode(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, database *db.TadDatabase, args []string) {
	if len(args) != 1 {
		reply(s, m, "usage: $owed @user")
		return
	}
	user, err := mentionedUser(m)
	if err != nil {
		reply(s, m, err.Error())
		return
	}

	owed, err := database.GetLeetcodeOwed(ctx, user.ID)
	if errors.Is(err, sql.ErrNoRows) {
		reply(s, m, user.DisplayName()+" doesn't owe any leetcode")
		return
	}
	if err != nil {
		log.Printf("owed: get owed for %s: %v", user.ID, err)
		reply(s, m, "something went wrong")
		return
	}
	reply(s, m, fmt.Sprintf("%s owes %d leetcode", user.DisplayName(), owed.OwedCount))
}

const leetcodeGraphQLURL = "https://leetcode.com/graphql"

const leetcodeQuery = `query($skip: Int, $limit: Int) {
	questionList(categorySlug: "", skip: $skip, limit: $limit, filters: {}) {
		total: totalNum
		questions: data { frontendQuestionId: questionFrontendId title titleSlug difficulty isPaidOnly }
	}
}`

// A random pick is paid roughly 1 in 6 times, so a few retries is plenty.
const maxLeetcodeAttempts = 5

func handleFetchRandomLeetcode(ctx context.Context, s *discordgo.Session, m *discordgo.MessageCreate, database *db.TadDatabase, args []string) {
	problem, err := randomLeetcodeProblem(ctx)
	if err != nil {
		log.Printf("leetcode: %v", err)
		reply(s, m, "couldn't fetch a leetcode problem, try again later")
		return
	}
	reply(s, m, formatLeetcodeProblem(problem))
}

// randomLeetcodeProblem picks a random free problem.
func randomLeetcodeProblem(ctx context.Context) (db.LeetcodeProblem, error) {
	first, err := fetchLeetcodeProblems(ctx, 0, 1)
	if err != nil {
		return db.LeetcodeProblem{}, err
	}
	if first.Total == 0 {
		return db.LeetcodeProblem{}, errors.New("no problems returned")
	}

	for range maxLeetcodeAttempts {
		list, err := fetchLeetcodeProblems(ctx, rand.IntN(first.Total), 1)
		if err != nil {
			return db.LeetcodeProblem{}, err
		}
		if len(list.Questions) > 0 && !list.Questions[0].PaidOnly {
			return list.Questions[0], nil
		}
	}
	return db.LeetcodeProblem{}, errors.New("only found paid problems")
}

func fetchLeetcodeProblems(ctx context.Context, skip, limit int) (db.LeetcodeProblemList, error) {
	body, err := json.Marshal(map[string]any{
		"query":     leetcodeQuery,
		"variables": map[string]int{"skip": skip, "limit": limit},
	})
	if err != nil {
		return db.LeetcodeProblemList{}, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, leetcodeGraphQLURL, bytes.NewReader(body))
	if err != nil {
		return db.LeetcodeProblemList{}, err
	}
	request.Header.Set("Content-Type", "application/json")
	// LeetCode rejects requests without a browser-like user agent and referer.
	request.Header.Set("User-Agent", "Mozilla/5.0")
	request.Header.Set("Referer", "https://leetcode.com")

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return db.LeetcodeProblemList{}, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return db.LeetcodeProblemList{}, fmt.Errorf("unexpected status %s", response.Status)
	}

	var result struct {
		Data struct {
			QuestionList db.LeetcodeProblemList `json:"questionList"`
		} `json:"data"`
	}
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return db.LeetcodeProblemList{}, fmt.Errorf("decode problems: %w", err)
	}
	return result.Data.QuestionList, nil
}

// formatLeetcodeProblem renders "#1 [Easy] Two Sum https://leetcode.com/problems/two-sum/".
func formatLeetcodeProblem(problem db.LeetcodeProblem) string {
	return fmt.Sprintf("#%s [%s] %s https://leetcode.com/problems/%s/",
		problem.Number, problem.Difficulty, problem.Title, problem.Slug)
}
