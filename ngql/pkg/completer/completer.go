/* Copyright (c) 2020 vesoft inc. All rights reserved.
 *
 * This source code is licensed under Apache 2.0 License.
 */

package completer

import (
	"strings"
)

// all keywords
var keywords = []string{
	"NODE",
	"TYPE",
	"TYPES",
	"TYPED",
	"FUNCTION",
	"GRAPH",
	"AS",
	"OR",
	"AND",
	"XOR",
	"USE",
	"SET",
	"FROM",
	"WHERE",
	"MATCH",
	"FILTER",
	"LET",
	"CALL",
	"INSERT",
	"YIELD",
	"RETURN",
	"DESCRIBE",
	"EDGE",
	"EDGES",
	"UPDATE",
	"UPSERT",
	"WHEN",
	"DELETE",
	"ALTER",
	"INDEX",
	"INDEXES",
	"REBUILD",
	"BOOL",
	"INT8",
	"INT16",
	"INT32",
	"INT64",
	"INT",
	"FLOAT",
	"DOUBLE",
	"STRING",
	"TIMESTAMP",
	"DATE",
	"DATETIME",
	"LIST",
	"RECORD",
	"UNION",
	"INTERSECT",
	"MINUS",
	"SHOW",
	"ADD",
	"CREATE",
	"DROP",
	"REMOVE",
	"CLOSE",
	"SESSION",
	"KILL",
	"QUERY",
	"COPY",
	"REPLACE",
	"IF",
	"NOT",
	"EXISTS",
	"WITH",
	"CHANGE",
	"GRANT",
	"REVOKE",
	"ON",
	"BY",
	"IN",
	"OF",
	"ORDER",
	"INGEST",
	"COMPACT",
	"FLUSH",
	"SUBMIT",
	"ASC",
	"DISTINCT",
	"BALANCE",
	"LIMIT",
	"OFFSET",
	"IS",
	"NULL",
	"EXPLAIN",
	"PROFILE",
	"FORMAT",
	"CASE",
	"MATCH",
	"SKIP",
	"SIGN",
	"HOSTS",
	"VALUES",
	"USER",
	"USERS",
	"PASSWORD",
	"ROLE",
	"ROLES",
	"GOD",
	"ADMIN",
	"DBA",
	"GUEST",
	"GROUP",
	"CHARSET",
	"COLLATE",
	"COLLATION",
	"ALL",
	"LEADER",
	"DATA",
	"SNAPSHOT",
	"SNAPSHOTS",
	"OFFLINE",
	"ACCOUNT",
	"JOBS",
	"JOB",
	"COUNT",
	"SUM",
	"AVG",
	"MAX",
	"MIN",
	"STD",
	"BIT_AND",
	"BIT_OR",
	"BIT_XOR",
	"PATH",
	"BIDIRECT",
	"STATUS",
	"FORCE",
	"PART",
	"PARTS",
	"DEFAULT",
	"HDFS",
	"CONFIGS",
	"PARAMS",
	"TTL_DURATION",
	"TTL_COL",
	"GRAPH",
	"META",
	"STORAGE",
	"SHORTEST",
	"CONTAINS",
	"TRUE",
	"FALSE",
	"THEN",
	"ELSE",
	"END",
	"STARTS",
	"ENDS",
	"WITH",
	"<-",
	"->",
	// ".",
	// ",",
	// ":",
	// ";",
	// "@",
	// "+",
	// "-",
	// "*",
	// "/",
	// "%",
	// "!",
	// "^",
	// "<",
	// "<=",
	// ">",
	// ">=",
	// "==",
	// "!=",
	// "||",
	// "&&",
	// "|",
	// "=",
	// "(",
	// ")",
	// "[",
	// "]",
}

var subCmds = map[string][]string{
	/* SHOW */
	"SHOW": []string{
		"HOSTS",
		"CHARSET",
		"COLLATION",
		"USERS",
		"CONFIGS",
		"CREATE",
		"INDEXES",
		"ROLES",
		"GRAPHS",
		"FUNCTIONS",
		"PLUGINS",
	},

	"CONFIGS": []string{
		"GRAPH",
		"STORAGE",
		"META",
	},

	"SHORTEST": []string{"PATH"},
	"ALL":      []string{"PATH"},
	"PATH":     []string{"FROM"},

	/* GROUP BY */
	"GROUP": []string{"BY"},

	/* ORDER BY */
	"ORDER": []string{"BY"},

	/* UNION */
	"UNION": []string{
		"DISTINCT",
		"ALL",
	},

	/* BALANCE */
	"BALANCE": []string{
		"LEADER",
		"DATA",
		"STOP",
	},

	/* DESCRIBE */
	"DESCRIBE": []string{
		"NODE",
		"EDGE",
		"GRAPH",
	},

	/* DDL */
	"CREATE": []string{
		"GRAPH",
		"USER",
	},
	"GRAPH": []string{
		"IF NOT EXISTS",
		"IF EXISTS",
	},
	"USER": []string{
		"IF NOT EXISTS",
		"IF EXISTS",
	},
	"INDEX": []string{
		"IF NOT EXISTS",
		"IF EXISTS",
	},
	"DROP": []string{
		"GRAPH",
		"USER",
	},

	/* DQL */
	"YIELD":  []string{"DISTINCT"},
	"STARTS": []string{"WITH"},
	"ENDS":   []string{"WITH"},

	/* DML */
	"INSERT": []string{
		"NODE",
		"EDGE",
	},

	/* something about user and role */
	"WITH":   []string{"PASSWORD"},
	"CHANGE": []string{"PASSWORD"},
	"GRANT":  []string{"ROLE"},
	"REVOKE": []string{"ROLE"},
}

func NewCompleter(line string, pos int) (head string, completions []string, tail string) {
	if len(line) < 1 {
		return
	}
	words := strings.Fields(line[:pos])
	if len(words) < 1 {
		return
	}
	lastWord := strings.ToUpper(words[len(words)-1])
	h := strings.LastIndex(line[:pos], " ")
	head = line[:h+1]
	tail = line[pos:]
	if line[pos-1] == ' ' { // find sub cmd
		if subs, ok := subCmds[lastWord]; ok {
			completions = append(completions, subs...)
		}
	} else {
		for _, k := range keywords {
			if strings.HasPrefix(k, lastWord) {
				completions = append(completions, k)
			}
		}
	}

	if len(completions) == 1 {
		completions[0] += " "
	}

	return
}
