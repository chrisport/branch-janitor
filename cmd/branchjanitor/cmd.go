package main

import (
	"github.com/chrisport/utils/exek"
	"strings"
	"fmt"
	"bufio"
	"os"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	branches := getBranches()
	if len(branches) == 0 {
		fmt.Println("Nothing to delete.")
		return
	}

	branchesToDelete := make([]string, 0)
	for i := range branches {
		r := ReadYesOrNo(reader, fmt.Sprintf("Delete %s\t[y/n] ", branches[i]))
		if r {
			branchesToDelete = append(branchesToDelete, branches[i])
		}
	}

	if len(branchesToDelete) == 0 {
		fmt.Println("Nothing to delete.")
		return
	}

	q := fmt.Sprintf("\n\nDelete [%s] from LOCAL:\t[y/n] ", strings.Join(branchesToDelete, ", "))
	yes := ReadYesOrNo(reader, q)
	if yes {
		fmt.Println(exek.Call("git branch -D " + strings.Join(branchesToDelete, " ")))
	}

	q = fmt.Sprintf("\n\nDelete [%s] from REMOTE:\t[y/n] ", strings.Join(branchesToDelete, ", "))
	yes = ReadYesOrNo(reader, q)
	if yes {
		fmt.Println(deleteRemote(branchesToDelete))
	}
}

func deleteRemote(branchesToDelete []string) string {
	origin := exek.Call("git branch -r")
	remoteBs := branchesToDelete[:0]
	for _, b := range branchesToDelete {
		if strings.Contains(origin, "origin/"+b) {
			remoteBs = append(remoteBs, b)
		}
	}
	fmt.Println("git push origin :" + strings.Join(remoteBs, " :"))
	return exek.Call("git push origin :" + strings.Join(remoteBs, " :"))
}

func ReadYesOrNo(reader *bufio.Reader, question string) bool {
	fmt.Print(question)
	t, err := reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	t = strings.ToLower(t)
	if t == "y\n" || t == "yes\n" {
		return true
	}

	if t == "n\n" || t == "no\n" {
		return false
	}

	fmt.Println("Please type y or n.")
	return ReadYesOrNo(reader, question)
}

func getBranches() []string {
	b := exek.Call("git branch -v")
	if len(b) == 0 {
		fmt.Println("Could not read branches, are you in a Git repository?")
		os.Exit(1)
	}
	lines := strings.Split(b, "\n")
	branches := make([]string, 0)
	for i := range lines {
		f := strings.Split(lines[i], " ")
		if f[0] != "*" {
			branches = append(branches, f[2])
		}
		fmt.Println(lines[i])
	}
	return branches
}
