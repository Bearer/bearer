package curio.detectors.test

import future.keywords

detect_at contains detections if {
  # input
  detections := curio.evaluator.node_detections(input, "objects")

  # detections := [{"hello": "world"}]
}

	# objectDetections, err := evaluator.NodeDetections(node, "objects")
	# if err != nil {
	# 	return nil, err
	# }

	# var result []*detectiontypes.Detection

	# for _, object := range objectDetections {
	# 	data := object.Data.(objects.Data)
	# 	if data.Name != "user" {
	# 		continue
	# 	}

	# 	for _, property := range data.Properties {
	# 		propertyData := property.Data.(properties.Data)
	# 		if propertyData.Name == "first_name" {
	# 			result = append(result, &detectiontypes.Detection{
	# 				ContextNode: node,
	# 				MatchNode:   property.MatchNode,
	# 				Data:        Data{Name: "Person Name"},
	# 			})
	# 		}
	# 	}
	# }

	# return result, nil
