// Copyright (c) 2023 Gitpod GmbH. All rights reserved.
// Licensed under the GNU Affero General Public License (AGPL).
// See License.AGPL.txt in the project root for license information.

package common

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/bufbuild/connect-go"
	v1 "github.com/gitpod-io/gitpod/components/public-api/go/experimental/v1"
)

func SshConnectToWs(ctx context.Context, workspaceID string, runDry bool) error {
	gitpod, err := GetGitpodClient(ctx)
	if err != nil {
		return err
	}

	workspace, err := gitpod.Workspaces.GetWorkspace(ctx, connect.NewRequest(&v1.GetWorkspaceRequest{WorkspaceId: workspaceID}))
	if err != nil {
		return err
	}

	wsInfo := workspace.Msg.GetResult()

	if wsInfo.Status.Instance.Status.Phase != v1.WorkspaceInstanceStatus_PHASE_RUNNING {
		return fmt.Errorf("cannot connect, workspace is not running")
	}

	token, err := gitpod.Workspaces.GetOwnerToken(ctx, connect.NewRequest(&v1.GetOwnerTokenRequest{WorkspaceId: workspaceID}))
	if err != nil {
		return err
	}

	ownerToken := token.Msg.Token

	host := strings.Replace(wsInfo.Status.Instance.Status.Url, wsInfo.WorkspaceId, wsInfo.WorkspaceId+".ssh", -1)
	host = strings.Replace(host, "https://", "", -1)

	if runDry {
		fmt.Println("ssh", fmt.Sprintf("%s#%s@%s", wsInfo.WorkspaceId, ownerToken, host), "-o", "StrictHostKeyChecking=no")
		return nil
	}

	slog.Debug("Connecting to", wsInfo.Description)

	command := exec.Command("ssh", fmt.Sprintf("%s#%s@%s", wsInfo.WorkspaceId, ownerToken, host), "-o", "StrictHostKeyChecking=no")

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		return err
	}

	return nil
}

type Response struct {
	OK      bool    `json:"ok"`
	Desktop Desktop `json:"desktop"`
}

type Desktop struct {
	Link     string `json:"link"`
	Label    string `json:"label"`
	ClientID string `json:"clientID"`
	Kind     string `json:"kind"`
}

func OpenWsInPreferredEditor(ctx context.Context, workspaceID string) error {
	gitpod, err := GetGitpodClient(ctx)
	if err != nil {
		return err
	}

	workspace, err := gitpod.Workspaces.GetWorkspace(ctx, connect.NewRequest(&v1.GetWorkspaceRequest{WorkspaceId: workspaceID}))
	if err != nil {
		return err
	}

	wsUrl, err := url.Parse(workspace.Msg.Result.Status.Instance.Status.Url)
	if err != nil {
		return err
	}

	wsHost := wsUrl.Host

	u := url.URL{
		Scheme: "https",
		Host:   wsHost,
		Path:   "_supervisor/v1/status/ide/wait/true",
	}

	fmt.Println((u.String()))

	resp, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var response Response
	if err := json.Unmarshal(body, &response); err != nil {
		return err
	}

	fmt.Printf("%+v\n", response)

	if response.OK {
		if response.Desktop.Link == "" {
			return fmt.Errorf("failed to open workspace in editor (no desktop editor)")
		}
		url := response.Desktop.Link
		var cmd *exec.Cmd
		switch os := runtime.GOOS; os {
		case "darwin":
			cmd = exec.Command("open", url)
		case "linux":
			cmd = exec.Command("xdg-open", url)
		case "windows":
			cmd = exec.Command("cmd", "/c", "start", url)
		default:
			panic("unsupported platform")
		}

		err := cmd.Start()
		if err != nil {
			return fmt.Errorf("failed to open workspace in editor: %w", err)
		}
	} else {
		return fmt.Errorf("failed to open workspace in editor (workspace not ready yet)")
	}

	return nil
}
