package main

import (
	"github.com/labstack/echo"
	"math"
	"net/http"
	"strings"
)

type ArcGISLOD struct {
	Level      int     `json:"level"`
	Resolution float64 `json:"resolution"`
	Scale      float64 `json:"scale"`
}

type ArcGISSpatialReference struct {
	Wkid uint16 `json:"wkid"`
}

type ArcGISExtent struct {
	Xmin             float32                `json:"xmin"`
	Ymin             float32                `json:"ymin"`
	Xmax             float32                `json:"xmax"`
	Ymax             float32                `json:"ymax"`
	SpatialReference ArcGISSpatialReference `json:"spatialReference"`
}

type ArcGISLayerStub struct {
	Id                uint8   `json:"id"`
	Name              string  `json:"name"`
	ParentLayerId     int16   `json:"parentLayerId"`
	DefaultVisibility bool    `json:"defaultVisibility"`
	SubLayerIds       []uint8 `json:"subLayerIds"`
	MinScale          float32 `json:"minScale"`
	MaxScale          float32 `json:"maxScale"`
}

type ArcGISLayer struct {
	Id                uint8             `json:"id"`
	Name              string            `json:"name"`
	Type              string            `json:"type"`
	Description       string            `json:"description"`
	GeometryType      string            `json:"geometryType"`
	CopyrightText     string            `json:"copyrightText"`
	ParentLayer       interface{}       `json:"parentLayer"`
	SubLayers         []ArcGISLayerStub `json:"subLayers"`
	MinScale          float32           `json:"minScale"`
	MaxScale          float32           `json:"maxScale"`
	DefaultVisibility bool              `json:"defaultVisibility"`
	Extent            ArcGISExtent      `json:"extent"`
	HasAttachments    bool              `json:"hasAttachments"`
	HtmlPopupType     string            `json:"htmlPopupType"`
	DrawingInfo       interface{}       `json:"drawingInfo"`
	DisplayField      string            `json:"displayField"`
	Fields            []interface{}     `json:"fields"`
	TypeIdField       string            `json:"typeIdField"`
	Types             string            `json:"types"`
	Relationships     []interface{}     `json:"relationships"`
	Capabilities      string            `json:"capabilities"`
}

var WebMercatorSR = ArcGISSpatialReference{Wkid: 3857}

func GetArcGISService(c echo.Context) error {
	id, err := getServiceOr404(c)
	if err != nil {
		return err
	}

	tileset := tilesets[id]
	imgFormat := tileset.format
	metadata := tileset.metadata
	name := toString(metadata["name"])
	description := toString(metadata["description"])
	attribution := toString(metadata["attribution"])

	// TODO: make sure that min and max zoom always populated
	scaleFactor := 156543.033928
	dpi := 96 // TODO: extract from the image instead
	minZoom := metadata["minzoom"].(int)
	maxZoom := metadata["maxzoom"].(int)
	var lods []ArcGISLOD
	for i := minZoom; i <= maxZoom; i++ {
		resolution := scaleFactor / math.Pow(2, float64(i))
		lods = append(lods, ArcGISLOD{
			Level:      i,
			Resolution: resolution,
			Scale:      float64(dpi) * 39.37 * resolution, // 39.37 in/m
		})
	}

	bounds := metadata["bounds"].([]float32) // TODO: make sure this is always present
	extent := ArcGISExtent{
		Xmin:             bounds[0],
		Ymin:             bounds[1],
		Xmax:             bounds[2],
		Ymax:             bounds[3],
		SpatialReference: WebMercatorSR,
	}

	tileInfo := map[string]interface{}{
		"rows": 256,
		"cols": 256,
		"dpi":  dpi,
		"origin": map[string]float32{
			"x": -20037508.342787,
			"y": 20037508.342787,
		},
		"spatialReference": WebMercatorSR,
		"lods":             lods,
	}

	documentInfo := map[string]string{
		"Title":    name,
		"Author":   attribution,
		"Comments": "",
		"Subject":  "",
		"Category": "",
		"Keywords": toString(metadata["tags"]),
		"Credits":  toString(metadata["credits"]),
	}

	out := map[string]interface{}{
		"currentVersion":            "10.4",
		"id":                        id,
		"name":                      name,
		"mapName":                   name,
		"capabilities":              "Map,TilesOnly",
		"description":               description,
		"serviceDescription":        description,
		"copyrightText":             attribution,
		"singleFusedMapCache":       true,
		"supportedImageFormatTypes": strings.ToUpper(imgFormat),
		"units":                     "esriMeters",
		"layers": []ArcGISLayerStub{
			ArcGISLayerStub{
				Id:                0,
				Name:              name,
				ParentLayerId:     -1,
				DefaultVisibility: true,
				SubLayerIds:       nil,
				MinScale:          0,
				MaxScale:          0,
			},
		},
		"tables":              []string{},
		"spatialReference":    WebMercatorSR,
		"tileInfo":            tileInfo,
		"documentInfo":        documentInfo,
		"initialExtent":       extent,
		"fullExtent":          extent,
		"exportTilesAllowed":  false,
		"maxExportTilesCount": 0,
		"resampling":          false,
	}

	return c.JSON(http.StatusOK, out)
}

func GetArcGISServiceLayers(c echo.Context) error {
	id, err := getServiceOr404(c)
	if err != nil {
		return err
	}

	tileset := tilesets[id]
	metadata := tileset.metadata

	bounds := metadata["bounds"].([]float32) // TODO: make sure this is always present
	extent := ArcGISExtent{
		Xmin:             bounds[0],
		Ymin:             bounds[1],
		Xmax:             bounds[2],
		Ymax:             bounds[3],
		SpatialReference: WebMercatorSR,
	}

	// for now, just create a placeholder root layer
	var layers [1]ArcGISLayer

	layers[0] = ArcGISLayer{
		Id:            0,
		ParentLayer:   nil,
		Name:          toString(metadata["name"]),
		Description:   toString(metadata["description"]),
		Extent:        extent,
		CopyrightText: toString(metadata["attribution"]),
		HtmlPopupType: "esriServerHTMLPopupTypeAsHTMLText",
	}

	out := map[string]interface{}{
		"layers": layers,
	}

	return c.JSON(http.StatusOK, out)
}

func GetArcGISServiceLegend(c echo.Context) error {
	id, err := getServiceOr404(c)
	if err != nil {
		return err
	}

	tileset := tilesets[id]
	metadata := tileset.metadata

	// TODO: pull the legend from ArcGIS specific metadata tables
	var elements [0]interface{}
	var layers [1]map[string]interface{}

	layers[0] = map[string]interface{}{
		"layerId":   0,
		"layerName": toString(metadata["name"]),
		"layerType": "",
		"minScale":  0,
		"maxScale":  0,
		"legend":    elements,
	}

	out := map[string]interface{}{
		"layers": layers,
	}

	return c.JSON(http.StatusOK, out)
}
