package vc_grouper

// the "garden" field lists some details about the kindoms available to the players
// "garden": [{
// 		"_id": 1,
// 		"block_x": 6,
// 		"block_y": 6,
// 		"unlock_block_x": 3,
// 		"unlock_block_y": 3,
// 		"bg_id": 1,
// 		"debris": 0,
// 		"castle_id": 7
// 	}, {
// 		"_id": 2,
// 		"block_x": 6,
// 		"block_y": 6,
// 		"unlock_block_x": 6,
// 		"unlock_block_y": 6,
// 		"bg_id": 2,
// 		"debris": 1,
// 		"castle_id": 66
// 	}]

//"garden_debris" lists information about clearing debris from your kingdom
// "garden_debris": [{
// 		"_id": 1,
// 		"garden_id": 2,
// 		"structure_id": 71,
// 		"x": 24,
// 		"y": 10,
// 		"level_cap": 1,
// 		"unlock_area_id": -1,
// 		"time": 1800,
// 		"coin": 8000,
// 		"iron": 8000,
// 		"ether": 8000,
// 		"cash": 0,
// 		"exp": 100
// 	}

// "structures" gives information about availability of for buildinds.
// The names of the structions in this list match those in the MsgBuildingName_en.strb file
// "structures": [{
// 		"_id": 1,
// 		"structure_type_id": 1,
// 		"max_lv": 10,
// 		"unlock_castle_id": 7,
// 		"unlock_castle_lv": -1,
// 		"unlock_area_id": -1,
// 		"base_num": 2,
// 		"size_x": 2,
// 		"size_y": 2,
// 		"order": 1000,
// 		"event_id": -1,
// 		"visitable": 0,
// 		"step": 0,
// 		"passable": 0,
// 		"connectable": 0,
// 		"enable": 1,
// 		"stockable": 1,
// 		"flag": 48,
// 		// 1 for kingdom 1, 2 for kingdom 2, 3 for both
// 		"garden_flag": 3
// 	}

// "event_structures" lists any structures available in the current event

//structure_level lists the level for the available structures
// "structure_level": [{
// 		"_id": 1,
// 		"structure_id": 28,
// 		"level": 1,
// 		"tex_id": 55,
// 		"level_cap": 1,
// 		"unlock_area_id": -1,
// 		"time": 0,
// 		"beginner_time": 0,
// 		"coin": 0,
// 		"iron": 0,
// 		"ether": 0,
// 		"cash": 500,
// 		"price": 0,
// 		"exp": 300
// 	}
