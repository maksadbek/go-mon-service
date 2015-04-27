function Canvas(map, latlng, msg, fillcolor, type, angle) {
 // Initialization
 this.setMap(map);
 this.map_ = map;
 var div = this.div_ = document.createElement('canvas');
	 div.style.cssText ='position: absolute;  left: 0; top: 0; white-space: nowrap; border: 0px solid blue; padding: 2px;  font-size:9pt';

	this.type_ =  type;
	
	if(this.type_ == "arrow")
	this.hw = 30
	else
	this.hw = 9;
	
	this.d = 8;
	div.width = this.hw;				  
	div.style.width = this.hw;				  
	div.id = 'can';
	div.height = this.hw;				  
	div.style.height = this.hw;				  
	this.ctx_ = this.div_.getContext('2d');
	this.cx =  (div.width)/2;
	this.cy =  (div.height)/2;
	
	this.angle_ = angle;
	this.latlng_ = latlng;
	//Canvas.infowindow_ = null;
	this.msg_ = '' + msg;
	this.fillcolor_ = ''+ fillcolor;
	var thisObj = this;
		
	
	//this.div_.onmouseover = function(event) 
	this.div_.onclick = function(event) 
	{
		r = 35/2;
		rx = event.offsetX - r;
		ry = event.offsetY - r;
		val = Math.pow(rx, 2) + Math.pow(ry, 2);

		if (Canvas.infowindow_)
		Canvas.infowindow_.close();
		
		Canvas.infowindow_ = new google.maps.InfoWindow({ position:thisObj.latlng_} );
		Canvas.infowindow_.setContent("<div class='infowin-content' style='height:205px'>" + thisObj.msg_ + "</div>");
		Canvas.infowindow_.open(thisObj.map_);
		Geocoding.getAddress(thisObj.latlng_, "span_");
			//var_dump(event);
			
		//document.getElementById("Memo").value = 'x=' + event.offsetX + ', y=' + event.offsetY + ', val=' + val + ', r=' + Math.pow(r, 2);	
	}
	
	this.div_.onmouseover = function(event) 
	{
		map.setOptions({draggableCursor:'pointer'});
	}

	this.div_.onmouseout = function(event) 
	{
		map.setOptions({draggableCursor:'default'});
	}
	
	/*this.div_.onmouseout = function(event) {
	
		r = 35/2;
		rx = event.offsetX - r;
		ry = event.offsetY - r;
		val = Math.pow(rx, 2) + Math.pow(ry, 2);
		
		if (thisObj.infowindow_)
		thisObj.infowindow_.close();
		
		document.getElementById("Memo").value = 'x=' + event.offsetX + ', y=' + event.offsetY + ', val=' + val + ', r=' + Math.pow(r, 2);	
		
	};
	
	
	/*
	tooltip =  new Label(
						this.map,
						latlng,
						info);
	*/
	
}
 

Canvas.prototype = new google.maps.OverlayView;
Canvas.infowindow_ = null;
Canvas.tooltip = null;




// Implement onAdd
Canvas.prototype.onAdd = function() {
	 
	
  //var pane = this.getPanes().overlayLayer;
  var pane = this.getPanes().overlayImage;
   pane.appendChild(this.div_);

	

 // Ensures the label is redrawn if the text or position is changed.
 var me = this;
    
 this.listeners_ = [
   google.maps.event.addListener(this, 'position_changed',
       function() { me.draw(); }),
   google.maps.event.addListener(this, 'text_changed',
       function() { me.draw(); })
 ];
};

// Implement onRemove
Canvas.prototype.onRemove = function() {
	
	//if (Canvas.infowindow_)
	//Canvas.infowindow_.close();

	if(this.div_ != null)
 	this.div_.parentNode.removeChild(this.div_);

 // Label is removed from the map, stop updating its position/text.
 for (var i = 0, I = this.listeners_.length; i < I; ++i) {
   google.maps.event.removeListener(this.listeners_[i]);
 }
};

Canvas.prototype.onMove = function(latlng, msg, fillcolor, type, angle){
	this.latlng_ = latlng;
	this.angle_ = angle;
	this.type_ =  type;
	this.msg_ = '' + msg;
	this.fillcolor_ = ''+ fillcolor;
	this.draw();	
};


// Implement draw
Canvas.prototype.draw = function() {
	
	var projection = this.getProjection();
	var pos =   new google.maps.LatLng(this.latlng_.lat(), this.latlng_.lng());
	position = projection.fromLatLngToDivPixel(pos);
	
	this.div_.width = this.div_.width;
	//var div = this.div_; 
	var x0 = (position.x - parseFloat(this.div_.width)/2 - 2 );
	var y0 = (position.y - parseFloat(this.div_.height)/2 - 2);
	
	this.div_.style.left = (x0) ;
	this.div_.style.top = (y0) ;
	this.div_.style.display = 'block';
	//var ctx = this.ctx_;

	this.ctx_.fillStyle = this.fillcolor_; // цвет фона
	this.ctx_.lineWidth = 0.4;
	this.ctx_.save();  
	this.ctx_.moveTo(0, 0);
	var cx = this.cx;
	var cy = this.cy;
	var d = this.d;

	if (this.type_ == "arrow")
	{
		this.ctx_.translate(cx, cy);
		this.ctx_.rotate(this.angle_* Math.PI / 180);	
		this.ctx_.translate(-cx, -cy);
		this.ctx_.beginPath();
		this.ctx_.moveTo(cx, cy);
		this.ctx_.lineTo(cx + d, cy + d);
		this.ctx_.lineTo(cx , cy - 2*d);
		this.ctx_.lineTo(cx - d, cy + d);
		this.ctx_.closePath();
		this.ctx_.fill();
		this.ctx_.stroke();
	}
	else
	{
		this.ctx_.beginPath();
		this.ctx_.arc(cx,cy,3,0,Math.PI*2,true); 
		this.ctx_.closePath();
		this.ctx_.fill();
		this.ctx_.stroke();
	}
	
};


	

