package utils

import (
	"time"

	"github.com/bwmarrin/snowflake"
)

func ToStringUuid() string {
	time.Sleep(1 * time.Millisecond)
	node, err := snowflake.NewNode(1)
	if err != nil {
		return ""
	}
	// Generate a snowflake ID.
	id := node.Generate()
	return id.String()
}
