package Util

import "github.com/GUAIK-ORG/go-snowflake/snowflake"

func GetGlobalId() int64 {
	s, err := snowflake.NewSnowflake(int64(0), int64(0))
	if err != nil {
		panic(err)
	}
	return s.NextVal()
}
