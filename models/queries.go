package datastore

var queries map[string]string = map[string]string{
	"usrTrackers": `SELECT 
                            login, 
                            fleet, 
                            cars
                        FROM max_users
                        WHERE login = ?
                        LIMIT 1`,
	"getCalibres": `
					SELECT 
						car_id,
						fleet_id,
						litr,
						volt
					FROM max_fuel_calibration`,
	"getTopLitres": `
				SELECT car_id , max( litr )
				FROM max_fuel_calibration
				WHERE 1
				GROUP BY car_id
				ORDER BY car_id 
				DESC
	`,
	"getTrackers": `select 
						u.id,
						u.fleet,
						u.imei,
						u.number,
						u.tracker_type,
						u.tracker_type_id,
						u.device_type_id,
						u.name,
						u.owner,
						u.active,
						u.additional,
						u.customization,
						u.group_id,
						u.detector_fuel_id,
						u.detector_motion_id,
						u.detector_dinamik_id,
						u.pid,
						u.installed_sensor,
						u.detector_agro_id,
						u.car_health,
						u.color,
						u.what_class,
						a.a_param_id 
					from
						max_units u
					LEFT JOIN additional_decode a
					 ON a.pr_name = 'fuelvolt' 
				  	 AND u.device_type_id = a.device_type
					 AND u.tracker_type_id = a.tracker_type
					 WHERE u.active = '1'`,
	"fleetTrackers": `
                      SELECT fleet, GROUP_CONCAT(id) cars
                      FROM max_units
                      WHERE active = '1' and fleet > 0
                      GROUP BY fleet
                      `,
	"checkUser": `
			SELECT *
			FROM max_users
			WHERE login = ?
			AND pass = ?
	`,
	`trackerGroups`: `select id, name, fllet_id from max_groups_units  order by fllet_id,name`,
}
