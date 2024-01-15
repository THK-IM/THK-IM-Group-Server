package test

import (
	"github.com/thk-im/thk-im-group-server/pkg/logic"
	"testing"
)

func TestGroupAvatar(t *testing.T) {
	g := logic.NewGroupAvatarGenerator("tmp", "test", "1.png")
	urls := []string{
		"tmp/1746103042618430147_1705138384708_0.png",
		"tmp/1746103042618430147_1705138384708_1.png",
		"tmp/1746103042618430147_1705138384708_1.png",
		"tmp/1746103042618430147_1705138384708_0.png",
		"tmp/1746103042618430147_1705138384708_1.png",
		"tmp/1746103042618430147_1705138384708_1.png",
		"tmp/1746103042618430147_1705138384708_0.png",
		"tmp/1746103042618430147_1705138384708_1.png",
		"tmp/1746103042618430147_1705138384708_1.png",
	}
	_, err := g.Compose(urls, 6)
	if err != nil {
		t.Failed()
	} else {
		t.Skip()
	}

}
