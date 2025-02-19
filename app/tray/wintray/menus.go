//go:build windows

package wintray

import (
	"fmt"
	"log/slog"
	"unsafe"

	"golang.org/x/sys/windows"
)

const (
	updatAvailableMenuID = 1
	updateMenuID         = updatAvailableMenuID + 1
	separatorMenuID      = updateMenuID + 1
	diagLogsMenuID       = separatorMenuID + 1
	diagSeparatorMenuID  = diagLogsMenuID + 1
	quitMenuID           = diagSeparatorMenuID + 1
)

func (t *winTray) initMenus() error {
	if err := t.addOrUpdateMenuItem(diagLogsMenuID, 0, diagLogsMenuTitle, false); err != nil {
		return fmt.Errorf("unable to create menu entries %w\n", err)
	}
	if err := t.addSeparatorMenuItem(diagSeparatorMenuID, 0); err != nil {
		return fmt.Errorf("unable to create menu entries %w", err)
	}
	if err := t.addOrUpdateMenuItem(quitMenuID, 0, quitMenuTitle, false); err != nil {
		return fmt.Errorf("unable to create menu entries %w\n", err)
	}
	return nil
}

func (t *winTray) UpdateAvailable(ver string) error {
	slog.Debug("updating menu and sending notification for new update")
	if err := t.addOrUpdateMenuItem(updatAvailableMenuID, 0, updateAvailableMenuTitle, true); err != nil {
		return fmt.Errorf("unable to create menu entries %w", err)
	}
	if err := t.addOrUpdateMenuItem(updateMenuID, 0, updateMenutTitle, false); err != nil {
		return fmt.Errorf("unable to create menu entries %w", err)
	}
	if err := t.addSeparatorMenuItem(separatorMenuID, 0); err != nil {
		return fmt.Errorf("unable to create menu entries %w", err)
	}
	iconFilePath, err := iconBytesToFilePath(wt.updateIcon)
	if err != nil {
		return fmt.Errorf("unable to write icon data to temp file: %w", err)
	}
	if err := wt.setIcon(iconFilePath); err != nil {
		return fmt.Errorf("unable to set icon: %w", err)
	}

	t.pendingUpdate = true
	// Now pop up the notification
	if !t.updateNotified {
		t.muNID.Lock()
		defer t.muNID.Unlock()
		copy(t.nid.InfoTitle[:], windows.StringToUTF16(updateTitle))
		copy(t.nid.Info[:], windows.StringToUTF16(fmt.Sprintf(updateMessage, ver)))
		t.nid.Flags |= NIF_INFO
		t.nid.Timeout = 10
		t.nid.Size = uint32(unsafe.Sizeof(*wt.nid))
		err = t.nid.modify()
		if err != nil {
			return err
		}
		t.updateNotified = true
	}
	return nil
}
