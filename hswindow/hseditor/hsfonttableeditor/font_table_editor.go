// Package hsfonttableeditor represents fontTableEditor's window
package hsfonttableeditor

import (
	"github.com/OpenDiablo2/dialog"

	g "github.com/ianling/giu"

	"github.com/OpenDiablo2/OpenDiablo2/d2common/d2fileformats/d2font"

	"github.com/OpenDiablo2/HellSpawner/hscommon"
	"github.com/OpenDiablo2/HellSpawner/hscommon/hsproject"
	"github.com/OpenDiablo2/HellSpawner/hswidget/fonttablewidget"
	"github.com/OpenDiablo2/HellSpawner/hswindow/hseditor"
)

const (
	mainWindowW, mainWindowH = 500, 400
)

const (
	removeItemButtonPath = "3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_delete.png"
	upItemButtonPath     = "3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_up.png"
	downItemButtonPath   = "3rdparty/iconpack-obsidian/Obsidian/actions/16/stock_down.png"
)

// static check, to ensure, if font table editor implemented editoWindow
var _ hscommon.EditorWindow = &FontTableEditor{}

// FontTableEditor represents font table editor
type FontTableEditor struct {
	*hseditor.Editor
	fontTable     *d2font.Font
	textureLoader *hscommon.TextureLoader
	textures
}

type textures struct {
	up,
	down,
	del *g.Texture
}

// Create creates a new font table editor
func Create(tl *hscommon.TextureLoader,
	pathEntry *hscommon.PathEntry,
	data *[]byte, x, y float32, project *hsproject.Project) (hscommon.EditorWindow, error) {
	table, err := d2font.Load(*data)
	if err != nil {
		return nil, err
	}

	result := &FontTableEditor{
		Editor:        hseditor.New(pathEntry, x, y, project),
		fontTable:     table,
		textureLoader: tl,
	}

	tl.CreateTextureFromFileAsync(removeItemButtonPath, func(texture *g.Texture) {
		result.textures.del = texture
	})

	tl.CreateTextureFromFileAsync(upItemButtonPath, func(texture *g.Texture) {
		result.textures.up = texture
	})

	tl.CreateTextureFromFileAsync(downItemButtonPath, func(texture *g.Texture) {
		result.textures.down = texture
	})

	return result, nil
}

// Build builds a font table editor's window
func (e *FontTableEditor) Build() {
	e.IsOpen(&e.Visible).Flags(g.WindowFlagsHorizontalScrollbar).
		Size(mainWindowW, mainWindowH).Layout(g.Layout{
		fonttablewidget.Create(e.textures.up, e.textures.down, e.textures.del, e.Path.GetUniqueID(), e.fontTable),
	})
}

// UpdateMainMenuLayout updates mainMenu layout's to it contain FontTableEditor's options
func (e *FontTableEditor) UpdateMainMenuLayout(l *g.Layout) {
	m := g.Menu("Font Table Editor").Layout(g.Layout{
		g.MenuItem("Add to project").OnClick(func() {}),
		g.MenuItem("Remove from project").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Import from file...").OnClick(func() {}),
		g.MenuItem("Export to file...").OnClick(func() {}),
		g.Separator(),
		g.MenuItem("Close").OnClick(func() {
			e.Cleanup()
		}),
	})

	*l = append(*l, m)
}

// GenerateSaveData generates data to be saved
func (e *FontTableEditor) GenerateSaveData() []byte {
	data := e.fontTable.Marshal()

	return data
}

// Save saves an editor
func (e *FontTableEditor) Save() {
	e.Editor.Save(e)
}

// Cleanup hides an editor
func (e *FontTableEditor) Cleanup() {
	if e.HasChanges(e) {
		if shouldSave := dialog.Message("There are unsaved changes to %s, save before closing this editor?",
			e.Path.FullPath).YesNo(); shouldSave {
			e.Save()
		}
	}

	e.Editor.Cleanup()
}
