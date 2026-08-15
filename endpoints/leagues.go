package endpoints

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"

	"github.com/AbdoAnss/go-fantasy-pl/api"
	"github.com/AbdoAnss/go-fantasy-pl/models"
)

const (
	classicLeagueEndpoint = "/leagues-classic/%d/standings/?page_standings=%d"
	h2hLeagueEndpoint     = "/leagues-h2h-matches/league/%d/"
	maxPageCache          = 3 // Only cache first 3 pages
)

// ErrLeagueNotFound is returned when the FPL API responds with HTTP 404,
// meaning the league ID does not exist (or is no longer accessible).
var ErrLeagueNotFound = errors.New("league not found")

// ErrInvalidH2HQuery is returned when the FPL API rejects the request with
// HTTP 400. In practice this happens for invalid query parameters, most
// commonly an event value outside the gameweek range (e.g. event=999).
type ErrInvalidH2HQuery struct {
	LeagueID    int
	QueryParams url.Values
	Detail      string
}

func (e *ErrInvalidH2HQuery) Error() string {
	return fmt.Sprintf("invalid query for H2H league %d (%s): the FPL API rejected the request (400)%s",
		e.LeagueID, encodeQuery(e.QueryParams), detailSuffix(e.Detail))
}

func encodeQuery(v url.Values) string {
	if len(v) == 0 {
		return "no query parameters"
	}
	return v.Encode()
}

func detailSuffix(detail string) string {
	if detail == "" {
		return ""
	}
	return ": " + detail
}

// LeagueService provides methods for fetching league standings and details,
// supporting both classic and head-to-head (H2H) leagues.
type LeagueService struct {
	client api.Client
}

// NewLeagueService creates a new instance of the LeagueService.
func NewLeagueService(client api.Client) *LeagueService {
	return &LeagueService{
		client: client,
	}
}

// GetClassicLeagueStandings returns the standings for a classic league by its unique ID.
// The page parameter allows for paginated access to large leagues (50 entries per page).
func (ls *LeagueService) GetClassicLeagueStandings(id, page int) (*models.ClassicLeague, error) {
	// Only cache first few pages to prevent memory bloat
	useCache := page <= maxPageCache

	if useCache {
		cacheKey := fmt.Sprintf("classic_league_%d_page_%d", id, page)
		var league models.ClassicLeague
		if sharedCache.Get(cacheKey, &league) {
			return &league, nil
		}
	}

	endpoint := fmt.Sprintf(classicLeagueEndpoint, id, page)
	resp, err := ls.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get league standings: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, fmt.Errorf("league with ID %d not found", id)
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var league models.ClassicLeague
	if err := json.Unmarshal(body, &league); err != nil {
		return nil, fmt.Errorf("failed to decode league data: %w", err)
	}

	if err := ls.validateLeague(&league); err != nil {
		return nil, err
	}

	if useCache {
		cacheKey := fmt.Sprintf("classic_league_%d_page_%d", id, page)
		if err := sharedCache.Set(cacheKey, &league, leagueCacheTTL); err != nil {
			return nil, fmt.Errorf("failed to cache league standings: %w", err)
		}
	}

	return &league, nil
}

func (ls *LeagueService) validateLeague(league *models.ClassicLeague) error {
	if league == nil {
		return fmt.Errorf("received nil league data")
	}
	if league.League.ID == 0 {
		return fmt.Errorf("invalid league ID")
	}
	return nil
}

// GetTotalPages calculates the total number of pages in a classic league.
func (ls *LeagueService) GetTotalPages(league *models.ClassicLeague) int {
	if league == nil || len(league.Standings.Results) == 0 {
		return 0
	}

	totalEntries := len(league.Standings.Results)
	if league.League.GetMaxEntries() > 0 {
		totalEntries = league.League.GetMaxEntries()
	}

	entriesPerPage := 50 // FPL default
	return (totalEntries + entriesPerPage - 1) / entriesPerPage
}

// H2HMatchesOption customises a GetH2HLeagueMatches request.
type H2HMatchesOption func(*h2hMatchesParams)

type h2hMatchesParams struct {
	event int
	page  int
}

// WithH2HEvent restricts results to a single gameweek (or knockout round
// when the value falls in the league's knockout range). Omit it to receive
// a paginated mixed feed spanning all played gameweeks. Values outside the
// valid gameweek range cause the FPL API to respond with HTTP 400, which
// surfaces as ErrInvalidH2HQuery.
func WithH2HEvent(event int) H2HMatchesOption {
	return func(p *h2hMatchesParams) { p.event = event }
}

// WithH2HPage requests a specific results page (1-indexed). Pagination
// applies both to the mixed feed and to event-filtered queries; note that
// has_next reflects the query used, so a filtered gameweek with few
// matches may report has_next=false where the mixed feed reports true.
func WithH2HPage(page int) H2HMatchesOption {
	return func(p *h2hMatchesParams) { p.page = page }
}

// GetH2HLeagueMatches returns head-to-head matches for a league.
//
// Without options it returns the first page of a paginated mixed feed that
// spans multiple gameweeks. With WithH2HEvent the feed is filtered to a
// single gameweek; for events in the knockout range the returned matches
// have is_knockout=true, a populated winner, and a knockout_name such as
// "Round 1". With WithH2HPage a specific page is fetched.
//
// Errors:
//   - ErrLeagueNotFound: HTTP 404, unknown league ID.
//   - *ErrInvalidH2HQuery: HTTP 400, e.g. an event value outside the
//     gameweek range (event=999).
func (ls *LeagueService) GetH2HLeagueMatches(leagueID int, opts ...H2HMatchesOption) (*models.H2HMatchesFeed, error) {
	params := h2hMatchesParams{}
	for _, opt := range opts {
		opt(&params)
	}

	query := url.Values{}
	if params.event != 0 {
		query.Set("event", strconv.Itoa(params.event))
	}
	if params.page != 0 {
		query.Set("page", strconv.Itoa(params.page))
	}

	endpoint := fmt.Sprintf(h2hLeagueEndpoint, leagueID)
	if len(query) > 0 {
		endpoint += "?" + query.Encode()
	}

	resp, err := ls.client.Get(endpoint)
	if err != nil {
		return nil, fmt.Errorf("failed to get H2H league matches: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
	case http.StatusNotFound:
		return nil, fmt.Errorf("%w: H2H league with ID %d", ErrLeagueNotFound, leagueID)
	case http.StatusBadRequest:
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			body = nil
		}
		return nil, &ErrInvalidH2HQuery{
			LeagueID:    leagueID,
			QueryParams: query,
			Detail:      parseDetail(body),
		}
	default:
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var feed models.H2HMatchesFeed
	if err := json.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("failed to decode H2H matches data: %w", err)
	}

	return &feed, nil
}

// parseDetail extracts a human-readable message from an error response body
// such as {"detail": "Not found."}.
func parseDetail(body []byte) string {
	var payload struct {
		Detail string `json:"detail"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return ""
	}
	return payload.Detail
}
