package blizzard

import (
	"context"
	"fmt"

	"github.com/FuzzyStatic/blizzard/v2/sc2gd"
)

// SC2LeagueData returns all SC2 league data from for seasonID, queue ID, team type, and league ID
func (c *Client) SC2LeagueData(ctx context.Context, seasonID int,
	queueID sc2gd.QueueID, teamType sc2gd.TeamType, leagueID sc2gd.LeagueID) (*sc2gd.League, []byte, error) {
	dat, b, err := c.getStructData(ctx,
		fmt.Sprintf("/data/sc2/league/%d/%d/%d/%d", seasonID, queueID, teamType, leagueID),
		"",
		&sc2gd.League{},
	)
	return dat.(*sc2gd.League), b, err
}

// SC2LadderData returns SC2 ladder for given division's ladderID.
// This API is undocumented by Blizzard, so it may be unstable.
func (c *Client) SC2LadderData(ctx context.Context, ladderID int) (*sc2gd.Ladder, []byte, error) {
	dat, b, err := c.getStructData(ctx,
		fmt.Sprintf("/data/sc2/ladder/%d", ladderID),
		"",
		&sc2gd.Ladder{},
	)
	return dat.(*sc2gd.Ladder), b, err
}
