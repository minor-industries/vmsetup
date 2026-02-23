package vmsetup

import (
	"os"
	"os/user"
	"strconv"
)

func chown(
	username string,
	groupname string,
	files ...string,
) error {
	u, err := user.Lookup(username)
	if err != nil {
		return err
	}

	g, err := user.LookupGroup(groupname)
	if err != nil {
		return err
	}

	uid, _ := strconv.Atoi(u.Uid)
	gid, _ := strconv.Atoi(g.Gid)

	for _, file := range files {
		if err := os.Chown(file, uid, gid); err != nil {
			return err
		}
	}

	return nil
}
