package datastore

var queries map[string]string = map[string]string{
	"usrTrackers": `SELECT 
                            login, 
                            fleet, 
                            cars
                        FROM max_users
                        WHERE login = ?
                        LIMIT 1`,
	"getTrackers": ` select 
                                 id,
                                 fleet,
                                 imei,
                                 number,
                                 tracker_type,
                                 tracker_type_id,
                                 device_type_id,
                                 name,
                                 owner,
                                 active,
                                 additional,
                                 customization,
                                 group_id,
                                 detector_fuel_id,
                                 detector_motion_id,
                                 detector_dinamik_id,
                                 pid,
                                 installed_sensor,
                                 detector_agro_id,
                                 car_health,
                                 color,
                                 what_class
                         from
                         max_units `,
	"fleetTrackers": `
                      SELECT fleet, GROUP_CONCAT(id) cars
                      FROM max_units
                      WHERE active = '1' and fleet > 0
                      GROUP BY fleet
                      `,
}
