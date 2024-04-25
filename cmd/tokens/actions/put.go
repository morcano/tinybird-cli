package tokens

import (
	"fmt"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"
	"log"
	"strings"
	"sync"
	"sync/atomic"
	"time"
	tokensapi "tinybird-cli/api/tokens"
)

var scope string

type job struct {
	token          tokensapi.Token
	completeChan   chan tokensapi.PutTokenParams
	failedCounter  *int32
	successCounter *int32
}

var PutToken = &cobra.Command{
	Use:   "put",
	Short: "Put a token",
	Long: `Modifies an Auth token. More than one scope can be sent per request, all of them will be added as Auth token scopes. 
	Everytime an Auth token scope is modified, it overrides the existing one(s).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		adminTok := cmd.Flag("admin-token").Value.String()
		client := tokensapi.NewClient(adminTok)

		fmt.Println("Getting tokens...")
		tokens, err := client.Get()
		if err != nil {
			log.Fatalf("Failed to get tokens: %v", err)
		}

		bulkUpdate(client, tokens.FilterByName("AccountToken").Tokens)

		return nil
	},
}

// bulkUpdate sends a PUT request to the TinyBird API using the provided auth tokens.
// It adds the specified scopes to the request URL and sets the necessary headers.
// If the request fails or returns a non-200 status code, it returns an error.
func bulkUpdate(client *tokensapi.Client, tokens []tokensapi.Token) {
	p := mpb.New(mpb.WithWidth(64))
	total := len(tokens)
	var successCount, failedCount int32
	var wg sync.WaitGroup

	numWorkers := 10
	jobs := make(chan job, total)
	completeChan := make(chan tokensapi.PutTokenParams, total)

	bar := p.AddBar(int64(total),
		mpb.PrependDecorators(
			decor.Name("Running... "),
			decor.Percentage(decor.WC{W: 5})),
		mpb.AppendDecorators(
			decor.Elapsed(decor.ET_STYLE_GO, decor.WC{W: 4}),
			decor.Any(func(s decor.Statistics) string {
				return fmt.Sprintf("Success: %d, Failed: %d", atomic.LoadInt32(&successCount), atomic.LoadInt32(&failedCount))
			}, decor.WC{W: 25}),
		),
	)

	// Spawn workers
	for w := 1; w <= numWorkers; w++ {
		go func(id int) {
			for j := range jobs {
				time.Sleep(500 * time.Millisecond) // Add delay here
				payload := tokensapi.PutTokenParams{
					Token: j.token.Token,
				}
				scopeList := strings.Split(scope, ",")
				for _, scope := range scopeList {
					payload.Scope = append(payload.Scope, scope+":accountId="+j.token.GetAccountIdFromName())
				}
				err := client.Put(payload)
				if err != nil {
					atomic.AddInt32(j.failedCounter, 1)
				} else {
					atomic.AddInt32(j.successCounter, 1)
				}
				j.completeChan <- payload
			}
		}(w)
	}

	// Submit jobs to workers
	for _, token := range tokens {
		jobs <- job{
			token:          token,
			completeChan:   completeChan,
			failedCounter:  &failedCount,
			successCounter: &successCount,
		}
		wg.Add(1)
	}
	close(jobs) // This signals to the workers that no more jobs are coming

	// Monitor completions
	go func() {
		for range completeChan {
			bar.Increment()
			wg.Done()
		}
	}()

	// Wait until all jobs are done
	wg.Wait()
	p.Wait()
}

func init() {
	PutToken.Flags().StringVarP(&scope, "scope", "s", "", "New scope(s) will override existing ones")
}
