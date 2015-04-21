package dbutil

/**
postgresify file provides the ability to convert queries done with question
marks into named queries with the proper query placeholders for postgres.
**/

import (
    "strconv"
    "os/exec"
)

var (
    postgresQueryConversions = make(map[string]string)
)

func QueryConvert(query, dbtype string) string {
    switch dbtype {
        case "postgres":
            return QueryToPostgres(query)
    }

    return query
}


// Converts all ? in a query to $n which is the postgres format
func QueryToPostgres(query string) string {

    // cacheing
    q := postgresQueryConversions[query]
    if q != "" {
        return q
    }

    output := ""
    position := 1

    for _, runeValue := range query {

        if runeValue == '?' {
            output += "$"
            output += strconv.Itoa(position)
            position += 1
            continue
        }

        output += string(runeValue)
    }

    return output
}

// finds the postgres binary on the system, isn't very robust in checking though
// should work on ubuntu variants and when postgres is on $PATH
func FindPostgres() string {

    // Start with which because we prefer a PATH version
    out := findPostgresWhich()

    if out != "" {
        return out
    }

    return findPostgresGrep()
}

// Find postgres using the lame grep method, works on Ubuntu (for now)
func findPostgresGrep() string {
    cmd := exec.Command("bash", "-c", "find /usr/lib/postgresql/ | sort -r | grep -m 1 /bin/postgres")
    out, err := cmd.CombinedOutput()

    if err != nil {
        return ""
    }

    return string(out)
}

// Finds postgres on $PATH
func findPostgresWhich() string {
    cmd := exec.Command("which", "postgres")
    out, err := cmd.CombinedOutput()

    if err != nil {
        return ""
    }

    return string(out)
}
