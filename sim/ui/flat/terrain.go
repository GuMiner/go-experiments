package flat

import (
	"go-experiments/sim/config"
	"go-experiments/sim/engine/terrain"
	"go-experiments/sim/ui/overlay"

	"github.com/go-gl/gl/v4.5-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type TerrainOverlay struct {
	textureId uint32
	overlay   *overlay.Overlay
}

func NewTerrainOverlay(textureId uint32) *TerrainOverlay {
	overlay := TerrainOverlay{
		textureId: textureId,
		overlay:   overlay.NewOverlay()}

	return &overlay
}

func (t *TerrainOverlay) GetOverlay() *overlay.Overlay {
	return t.overlay
}

type TerrainOverlayManager struct {
	TerrainOverlays map[int]map[int]*TerrainOverlay
}

func (t *TerrainOverlayManager) GetOrAddTerrainOverlay(x, y int) (*TerrainOverlay, bool) {
	isNew := false
	if _, ok := t.TerrainOverlays[x]; !ok {
		t.TerrainOverlays[x] = make(map[int]*TerrainOverlay)
	}

	if _, ok := t.TerrainOverlays[x][y]; !ok {
		var textureId uint32
		gl.GenTextures(1, &textureId)
		t.TerrainOverlays[x][y] = NewTerrainOverlay(textureId)
		isNew = true
	}

	return t.TerrainOverlays[x][y], isNew
}

func NewTerrainOverlayManager() *TerrainOverlayManager {
	manager := TerrainOverlayManager{
		TerrainOverlays: make(map[int]map[int]*TerrainOverlay)}

	return &manager
}

func (t *TerrainOverlayManager) Delete() {
	for _, value := range t.TerrainOverlays {
		for _, overlay := range value {
			gl.DeleteTextures(1, &overlay.textureId)
		}
	}
}

func (t *TerrainOverlay) UpdateCameraOffset(x, y int, camera *Camera) {
	offset := camera.GetRegionOffset(x, y)
	scale := camera.GetRegionScale()

	t.overlay.UpdateLocation(offset, scale, 1.0)
}

func (t *TerrainOverlay) SetTerrain(texels [][]terrain.TerrainTexel) {
	regionSize := len(texels[0])
	byteTerrain := make([]uint8, regionSize*regionSize*4)
	for i := 0; i < regionSize; i++ {
		for j := 0; j < regionSize; j++ {
			height := texels[i][j].Height

			color, percent := GetTerrainColor(height)
			byteTerrain[(i+j*regionSize)*4] = uint8(color.X() * percent)
			byteTerrain[(i+j*regionSize)*4+1] = uint8(color.Y() * percent)
			byteTerrain[(i+j*regionSize)*4+2] = uint8(color.Z() * percent)
			byteTerrain[(i+j*regionSize)*4+3] = 1.0
		}
	}

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, t.textureId)
	gl.TexStorage2D(gl.TEXTURE_2D, 1, gl.RGBA8, int32(regionSize), int32(regionSize))
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)

	gl.TexSubImage2D(gl.TEXTURE_2D, 0,
		0, 0, int32(regionSize), int32(regionSize),
		gl.RGBA, gl.UNSIGNED_BYTE, gl.Ptr(byteTerrain))

	t.overlay.UpdateTexture(t.textureId)
}

// Given a height, returns the terrain color and percentage within that level
func GetTerrainColor(height float32) (mgl32.Vec3, float32) {
	terrainType, percent := terrain.GetTerrainType(height)

	switch terrainType {
	case terrain.Water:
		return config.Config.Ui.TerrainUi.WaterColor.ToVec3(), percent
	case terrain.Sand:
		return config.Config.Ui.TerrainUi.SandColor.ToVec3(), percent
	case terrain.Grass:
		return config.Config.Ui.TerrainUi.GrassColor.ToVec3(), percent
	case terrain.Hills:
		return config.Config.Ui.TerrainUi.HillColor.ToVec3(), percent
	case terrain.Rocks:
		return config.Config.Ui.TerrainUi.RockColor.ToVec3(), percent
	default:
		return config.Config.Ui.TerrainUi.SnowColor.ToVec3(), percent
	}
}
