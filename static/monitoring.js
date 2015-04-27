//Monitoring.js
var UserBrowser = 1; 
//vqatincha tek uchun...
function Monitoring(map)
{
	
	this.map = map;
	objdirection:false;
	objname:false;
	smoothly:false;
	smooth_type:2;
	smooth_val:5;
	
	var  pois = [];
	var bsics = {};
	var zoomall = false;
	var bounds = new google.maps.LatLngBounds();
	var tmpcars = [];
	var geo_anims = {};
		
	function addPOI(newp)
	{
		var ind  = findById(newp.id); 
		
		if(ind == -1)
		{			
			pois.push(new MPoints(map, newp));
			ind = pois.length - 1;
		}
		else
		{
			
			pois[ind].todo = newp.action;
			if(newp.action != '-1')
			{
				pois[ind].setNewParams(newp);
			}
		}
			
		pois[ind].execute();
		//show in geozone alerts.
		if (geo_anims[newp.id] != null) 
		{
			hideShowAlert(newp.id, geo_anims[newp.id].show, newp.action);
		}
	};
	
	function var_dump(o)
	{	
		alert(print_r(o));
	}
	
	function print_r(o)
	{
	
		var r = '';
		for(var i in o) r += i + ': ' + o[i] + '\n';
		return r;
	}
	
	
	
	function addBSIC(cars)
	{
				
	
		var image = new google.maps.MarkerImage('/i/cell_tower/cell_tower_1.png');
		color = getColorbyId(cars.id);
				
		// on marker move and cell id move
		
		var car_id = findCarId(cars.id);		
				
		if (car_id != -1)
		{
		
			
			for(var i = 0;  i < tmpcars[car_id].cell_location.length; i++)
			{
				
				cell_id = -1;
				cid = tmpcars[car_id].cell_location[i].cellid;
			
				if (cars.action != '-1')
				{
				
					for(var j = 0, t = cars.cell_location.length; j < t; j++)
					{
						if (cars.cell_location[j].cellid == cid)
						{
							cell_id = j;
						}
					}
				
				}
			
				if (cell_id == -1)
				{
				
					if (bsics[cid].units.length == 1)
					{
						bsics[cid].obj.setMap(null);
						bsics[cid].marker.setMap(null);
						bsics[cid].poly[0].setMap(null);
						bsics[cid].poly.pop();
						bsics[cid].units.pop();
						bsics[cid].cellid  = 0;
					}
					else
					{

						id = bsics[cid].units.indexOf(cars.id*1);
						if (id != -1)
						{
						bsics[cid].units.splice(id, 1);
						bsics[cid].poly[id].setMap(null);
						bsics[cid].poly.splice(id, 1);
						}
					}
					
				}
				
				
				//var_dump(cars);
				
				// remove objects 
				if (cars.action == '-1')
				{
				
					id = bsics[cid].units.indexOf(cars.id*1);
					 //id = findUnitId(cid, cars.id);
					
					//alert("unit_id:" + cars.id + " id:" + id);
					//var_dump(bsics[cid].units);
					
					if (id != -1)
					{
						bsics[cid].units.splice(id, 1);
						bsics[cid].poly[id].setMap(null);
						bsics[cid].poly.splice(id, 1);
					
					}
					if (bsics[cid].units.length == 0)
					{
						bsics[cid].obj.setMap(null);
						bsics[cid].marker.setMap(null);
						bsics[cid].cellid  = 0;
					
					}
						
				}
				
			}
			
		
		}
		
		
	
		if (cars.action == '-1') 
		{
			
			return;
		}
				
			
		for(var i = 0, t = cars.cell_location.length; i < t; i++)
		{
			
			// Если из базы не может найти для это ID нужный Lon Lat то пропускает цикл не вставляет маркер
			if( (cars.cell_location[i].lat == null) || (cars.cell_location[i].lon == null) )
				continue;
			
			var latlng = new google.maps.LatLng(cars.cell_location[i].lat, cars.cell_location[i].lon);
			
			var id = findCellID(cars.cell_location[i].cellid);
			
			
				
			coords = [
						new google.maps.LatLng(cars.latitude, cars.longitude), 
						latlng
				 ];
				 
				 
			poly = new google.maps.Polyline({
							path: coords,
							strokeColor:color,
							strokeOpacity: 0.2,
							strokeWeight: 2,
							geodesic: true							
							
					});
				
			//  CellID is not exists...		
			if (id == -1)
			{			
			
				options = {
					  strokeOpacity: 0.8,
					  strokeWeight: 2,
					  fillOpacity: 0.5,
					  strokeColor: color,
					  fillColor: color,
					  center: latlng,
					  radius: 60
				};
							
				circle = new google.maps.Circle(options);
				
				m =  new google.maps.Marker({
					icon:image,
					position:latlng
				});
				
				
				bsics[cars.cell_location[i].cellid] = {
					latlng: latlng,
					cellid:  cars.cell_location[i].cellid,
					bsic: cars.cell_location[i].bsic,
					rssi: cars.cell_location[i].rssi,
				    obj:circle,
					poly:[],
					units:[],
					marker:m
				};
			
				bsics[cars.cell_location[i].cellid].units.push(cars.id*1);
				bsics[cars.cell_location[i].cellid].poly.push(poly);
				
				bsics[cars.cell_location[i].cellid].obj.setMap(map);
				bsics[cars.cell_location[i].cellid].marker.setMap(map);
				poly.setMap(map);
				//break;				
				continue;
			}

			// if cellID exists add  unit_id and polyline
			index = bsics[cars.cell_location[i].cellid].units.indexOf(cars.id*1);
			if (index < 0)
			{
				bsics[cars.cell_location[i].cellid].units.push(cars.id*1);
				bsics[cars.cell_location[i].cellid].poly.push(poly);
				poly.setMap(map);
				continue;
			}
			
			// on new position move..
			bsics[cars.cell_location[i].cellid].poly[index].setPath(coords);
						
		}
		
	};
	

	
	
	function show_props(obj, obj_name)
	{
		var result = "";
		for (var i in obj)
		{
			result += obj_name + "." + i + " = " + obj[i] + "\n";
		}
		return result;
	};
	

	
	this.obj_array = function(cars, zoom)
	{		
		bounds = new google.maps.LatLngBounds();
					
		for(var i = 0, t = cars.length; i < t; i++)
		{
			if (cars[i]!=undefined){
				var obj = {
					id: cars[i].id, 
					action: cars[i].action,
					sat:cars[i].sat
				};
			
				if(obj.action != '-1')
				{
					obj.car_name = cars[i].car_name;
                    obj.owner = cars[i].owner;
					obj.formatted_time = cars[i].formatted_time;
					obj.direction = cars[i].direction;
					obj.speed = cars[i].speed;
					obj.addparams = cars[i].addparams;
					obj.sat = cars[i].sat;
					obj.latitude = cars[i].latitude;
					obj.longitude = cars[i].longitude;
				
				}
				
	  			addPOI(obj);
		
				if (obj.sat == 100)
					addBSIC(cars[i]);
					
				if (cars[i].action != '-1')
				{
					tmpcars = cars;
					bounds.extend(new google.maps.LatLng(cars[i].latitude, cars[i].longitude));
				}
			}
		}
		
		if ( (zoom && cars.length>0) || zoomall)
		{
			if (!bounds.isEmpty())
			{
				//var_dump(bounds);
				this.map.fitBounds(bounds);
			}
			else 
			{
				if(LongitudeJs==undefined)
				{
					LatitudeJs = 41.328993414841975;
					LongitudeJs = 69.24995110359192;
			    }
			    
			    if(this.map!=undefined)
			    {
					this.map.setCenter(new google.maps.LatLng(LatitudeJs, LongitudeJs));
					this.map.setZoom(ZoomJs);
			    }
				
			}
		}
		
		
		
	};
	
		//geozone alerts on/off animation	 
	this.geozoneAlert = function(id, bool){
		hideShowAlert(id, bool, '0');		
	}


	function hideShowAlert(id, bool, action){
		
		//off animation
		if (!bool || action == '-1')
		{
			if (geo_anims[id] != null)
				{
					if (geo_anims[id].obj != null){
						geo_anims[id].obj.setMap(null);
						geo_anims[id].obj = null;
					}
					geo_anims[id].show = bool;
					

				}

			return;			
		}

		//create marker animation and update postion
		if (geo_anims[id] == null)
			geo_anims[id] = {
				obj:null
			};
			geo_anims[id].show = bool;
		
		i = findById(id);
		
		if (i >= 0 && action == '0') action = pois[i].todo;
		if (i < 0 || action == '-1') return;

		if (geo_anims[id].obj != null )
		{
			geo_anims[id].obj.setPosition(pois[i].coordinates);
			return;
		}

		var image = new google.maps.MarkerImage(
				'/i/monitoring/transparent.png',
				null, // size
				null, // origin
				new google.maps.Point( 8, 8 ), // anchor (move to center of marker)
				new google.maps.Size( 17, 17 ) 
			);	
		
		 marker = new google.maps.Marker({
			icon:image,
			flat:true,
			optimized:false,
			title:'anim',
			position:pois[i].coordinates
		});

		geo_anims[id].obj = marker;
		geo_anims[id].obj.setMap(map);

	};
	//geozone alerts on/off animation....	 

	
	this.hideObjects = function()
	{
		for(var i = 0, t = pois.length; i < t; i++)
		{
			pois[i].todo = -1;
			pois[i].execute();
		}
		pois.length = 0;
		
		for (var cid in bsics)
		{
			bsics[cid].obj.setMap(null);
			bsics[cid].marker.setMap(null);			
			for (var j = 0; j < bsics[cid].poly.length; j++)
				bsics[cid].poly[j].setMap(null);
		}
		
		bsics  = {};
		//clear  all geo animation
		for(var i in geo_anims)
		{
			if (geo_anims[i].obj != null)
			{
				geo_anims[i].obj.setMap(null);					
				geo_anims[i].obj = null;
			}
		}
		geo_anims = {};
		
	};
	
	
	
	this.showName = function(bool)
	{
		Monitoring.objname = bool;
		for(var i = 0, t = pois.length; i < t; i++)
		{
			pois[i].objshowname();
		}
	};
	
	this.setCenterObj = function(id)
	{
		var i = findById(id);
		if (i != -1)
		{
			if (pois[i].todo != -1)
			this.map.setCenter(pois[i].coordinates, 16);
		}
	};
		
	function removeEventListener()
	{
		//this.map.removeEventListener(MapEvent.FLY_TO_DONE, onFlyDone);	
	}
	
	this.zoomAll = function(t)
	{
		zoomall = t;
	};
		
	
	this.setSmoothly = function(smooth, type, val)
	{
		Monitoring.smoothly    = smooth;
		Monitoring.smooth_type = type;
		Monitoring.smooth_val  = val;
	
		for(var i = 0;  i < pois.length;  i++)
		{
			pois[i].removeSmoothLine();
		}
		
	};

		
	function obj_direction(bool)
	{
		Monitoring.objdirection = bool;
		for(var i = 0, t = pois.length; i < t; i++)
		{
			pois[i].options = pois[i].objOptions(pois[i].sprite, pois[i].direction);
			pois[i].obj.setOptions(pois[i].options);
		}
	};
	 
	
		
		
	function findById(id)
	{
		
		for(var i = 0;  i < pois.length;  i++)
		{
			if (pois[i].id == id) return i;
		}
			return -1;
			
	};	
	
	function findCarId(id)
	{
		for(var i = 0;  i < tmpcars.length;  i++)
		{
			if (tmpcars[i].id == id) return i;
		}
			return -1;
	};
	
	
	function findCellID(id)
	{
		 for (var i in bsics) {
			if (bsics[i].cellid == id)
			return i;
		 }
		 return -1;
	};
	
	function findUnitId(cid, id)
	{
		for (var  i = 0; i < bsics[cid].units.length; i++)
		{
			if (bsics[cid].units[i]  == id)
			return i;
		}
		
		return -1;
	};
	
	
	
}



function MPoints(map, arr_car)
{
	this.map = map;
	this.id = arr_car.id;
	this.todo = arr_car.action;
	var created = false;
	var obj = null;
	this.coordinates = 0;
	this.lastcoordinates = 0;
	this.htmlInfo = '';
	var infoOpt = arr_car;
	var thisOjb = this;
	//var infowindow;
	var objLabel = null;
	var m_icon = getIcon(infoOpt.id);
	var smooth_val = 0;
	//var m_icon = '0'; 
	
	//-------vaqtincha tekshirish uchun
	var image = new google.maps.MarkerImage(m_icon);
	day = new Date();
	var time_period = day.getTime();
	
	//smoothly move vars
	var timerHandle = null;
	var eol = 0;
	var polyline = new google.maps.Polyline({
	path: [],
	strokeColor: '#FF0000',
	strokeWeight: 3,
	strokeOpacity: 0.4
    });
	var line = new google.maps.Polyline({
	path: [],
	strokeColor: '#FF0000',
	strokeWeight: 5,
	strokeOpacity: 0.7,
	map:thisOjb.map
    });
	
	line.setOptions({strokeColor:getColorbyId(infoOpt.id)});
	

	
this.setNewParams = function(arr_car)
{
	this.coordinates = new google.maps.LatLng(arr_car.latitude, arr_car.longitude);
	this.carname = arr_car.car_name;
	this.direction = arr_car.direction;
	this.sat = arr_car.sat;
	this.addparams = arr_car.addparams;
	infoOpt  = arr_car;
	this.htmlInfo = getInfoOpt(infoOpt);
	
	if (arr_car.sat == 100)
	{
		cell_location = arr_car.cell_location;
	}
	if (this.lastcoordinates != 0 && Monitoring.smoothly)
	{
	polyline.getPath().clear();
	polyline.getPath().push(this.lastcoordinates);
	polyline.getPath().push(this.coordinates);
	}
	this.lastcoordinates = this.coordinates;
	
	
};

this.setNewParams(arr_car);


function hidePoi()
{
	if(created)
	{
		if(UserBrowser==1)
		{	
				if (m_icon == '0')
				{			
					obj.onRemove();
					obj = null;
					line.getPath().clear();
					created = false;
					
				}
				else
				{
					obj.setMap(null);
					created = false;	
					
				}
		} else {
			// IE 
			obj.setMap(null);
			created = false;	
		}		
	}
}
		


function newPoiPos()
{
	
	//htmlInfo = getInfoOpt(infoOpt);
	if (m_icon == '0')
	{
		console.log(m_icon);
		if (obj != null)
		{
			
			obj.onMove(thisOjb.coordinates, thisOjb.htmlInfo, getColorbyId(infoOpt.id), "arrow", thisOjb.direction);
		}
	}
	else
	{
		if (created)
		{
			obj.setPosition(thisOjb.coordinates);
			m_icon = getIcon(infoOpt.id);
			image = new google.maps.MarkerImage(m_icon);
			obj.setIcon(image);
		}
		else
		{
			obj =  new google.maps.Marker({
				icon:image,
				position:thisOjb.coordinates
			});
			
			obj.setMap(thisOjb.map);	
			 google.maps.event.addListener(obj, 'click',
       			function() {
					
					var infowindow_ = new google.maps.InfoWindow({ position:thisOjb.coordinates} );
					infowindow_.setContent("<div class='infowin-content' style='height:205px'>" + thisOjb.htmlInfo + "</div>");
					infowindow_.open(thisOjb.map);
					Geocoding.getAddress(thisOjb.coordinates, "span_");
       			});
       	
       			created = true;		
		}
	}
}
	
		
function destroy()
{
	hidePoi();
}




function addNewPoi()
{

	
	if(!created)
	{
		//var htmlInfo = getInfoOpt(infoOpt);
		if(UserBrowser==1)
		{
			if (m_icon == '0')
			{
				if (obj == null)
				obj = new Canvas(thisOjb.map, thisOjb.coordinates, thisOjb.htmlInfo, getColorbyId(infoOpt.id), "arrow", thisOjb.direction);
			}
			else
			{
				obj =  new google.maps.Marker({
					icon:image,
					position:thisOjb.coordinates,
					title: 'I might be here',
					visible: true
				});
				
				obj.setMap(thisOjb.map);	
				 google.maps.event.addListener(obj, 'click',
	       			function() {
	       					var infowindow_ = new google.maps.InfoWindow({ position:thisOjb.coordinates} );
							infowindow_.setContent("<div class='infowin-content' style='height:205px'>" + thisOjb.htmlInfo + "</div>");
							infowindow_.open(thisOjb.map);
							Geocoding.getAddress(thisOjb.coordinates, "span_");
						
						
						
	       			});
			}
		}else{
			// елси IE
			
			obj =  new google.maps.Marker({
					icon:'/i/monitoring/gprs-online.png',					
					position:thisOjb.coordinates,
					title: 'I might be here',
					flat:true
			});
			
			obj.setMap(thisOjb.map);	
			 google.maps.event.addListener(obj, 'click',
       			function() {
       					
       					var infowindow_ = new google.maps.InfoWindow({ position:thisOjb.coordinates} );
       					
						infowindow_.setContent("<div class='infowin-content' style='height:205px'>" + thisOjb.htmlInfo + "</div>");
						
						infowindow_.open(thisOjb.map);
						
       			});
			
		}
		
		
		created = true;

	}
	else
	{	
		newPoiPos();	
	}
			
}


function nextStep(d){
	if (d>eol) {
	clearTimeout(timerHandle);
	thisOjb.coordinates = polyline.getPath().getAt(1);
	line.getPath().push(thisOjb.coordinates);
	//thisOjb.direction = this.direction;
	removeVertex();
	newPoiPos();
	return;
	}
	moveObject(d);
	d += 50;
	timerHandle =  setTimeout(function() { nextStep(d); }, 100);
}


function moveObject(d){
	removeVertex();
	var p = polyline.GetPointAtDistance(d);
	thisOjb.coordinates = p;
	newPoiPos();
	line.getPath().push(p);
}

function removeVertex(){
	if (Monitoring.smooth_type == 2 || Monitoring.smooth_type == 1)
	{
		if (line.getPath().getLength() > 0 && smooth_val >= Monitoring.smooth_val) 
		line.getPath().removeAt(0);
	}
}

function startMove() {
	day = new Date();
	eol = polyline.Distance();
	//thisOjb.direction = polyline.Bearing(0, 1);
	if (smooth_val < Monitoring.smooth_val && eol >= 10)
	{
		if (Monitoring.smooth_type == 1) 
		{
			smooth_val = smooth_val +  (day.getTime() - time_period)/1000;
			time_period = day.getTime();
		}
		else
		if (Monitoring.smooth_type == 2)
			smooth_val++;	
	}
	
	nextStep(1);
}

this.removeSmoothLine = function()
{
	line.getPath().clear();
	smooth_val = 0;
};

function getInfoOpt(arr)
{
	
	var html = 	
	'<div align="left"><div style="margin-top:3px; font-size:12px; font-family: Tahoma;"> номер №: <b>' 	+ arr.id + '</b></div>' + 
	'<div style="margin-top:3px; font-size:12px; font-family: Tahoma;">name:<b> ' 	+ arr.car_name + '</b></div>' 
		
	return html;	
}	

this.execute = function()
{
	if(this.todo == '1')
	{
		if (Monitoring.smoothly && polyline.getPath().getLength() > 1)
		startMove();
		else
		newPoiPos();
	}
	else if(this.todo == '2')
	{
		addNewPoi();
	}
	else if(this.todo == '-1')
	{
		hidePoi();
	}
	
	this.objshowname();
};

this.objshowname = function()
{
	if (objLabel != null)	
	{
		objLabel.onRemove();
		objLabel = null;
	} 
		
	if (Monitoring.objname)
	{
		if (this.todo != '-1')
		{
			 if (objLabel != null)	
			  objLabel.onRemove();
			  
			  objLabel = new Label(
				this.map,
			this.coordinates, infoOpt.car_name );
	 			
		}
	}
	
	
	
};



function getIcon(id)
{
	for(var i = 0, t = MonitoringTools.arr_tools.length; i < t; i++)
	{
	  if (id == MonitoringTools.arr_tools[i].id)
	  {	  	
	  	if(MonitoringTools.arr_tools[i].icon_path == '')
			return '0';
		else
			return MonitoringTools.arr_tools[i].icon_path;
	  }
	}
	
	return '0';
}

	
}


function getColorbyId(id)
{
	for(var i = 0, t = MonitoringTools.arr_tools.length; i < t; i++)
	{
	  if (id == MonitoringTools.arr_tools[i].id)
	  {
		return MonitoringTools.arr_tools[i].fillcolor;
	  }
	}
	
	return '#D3005F';
}

	
	
	
	

