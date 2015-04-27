// Define the overlay, derived from google.maps.OverlayView



function Canvas2(map, latlng, msg, fillcolor, type, angle) {
 // Initialization
 //this.setMap(map);
 this.map_ = map;
 //this.setValues(opt_options);
  var div = this.div_ = document.createElement('canvas');
	 div.style.cssText ='position: absolute;  left: 0; top: 0; white-space: nowrap; border: 0px solid blue; padding: 2px;  font-size:9pt;';

	this.type_ =  type;
	
	if(this.type_ == "arrow")
	this.hw = 30
	else
	this.hw = 9;
	
	this.d = 4;
	div.width=this.hw;				  
	div.style.width = this.hw;				  
	div.id='can';
	div.height=this.hw;				  
	div.style.height=this.hw;				  
	this.ctx_=this.div_.getContext('2d');
	this.cx =(div.width)/2;
	this.cy =(div.height)/2;
	
	this.angle_ = angle;
	this.latlng_ = latlng;
	//Canvas.infowindow_ = null;
	this.msg_ = '' + msg;
	this.fillcolor_ = ''+ fillcolor;
	var thisObj = this;
		
	
	this.show = function()
	{
		this.setMap(map);
	}
	//this.div_.onmouseover = function(event) 
	this.div_.onclick = function(event) 
	{
/*
		r = 35/2;
		rx = event.offsetX - r;
		ry = event.offsetY - r;
		val = Math.pow(rx, 2) + Math.pow(ry, 2);

		if (Canvas2.infowindow_)
		Canvas2.infowindow_.close();
		
		Canvas2.infowindow_ = new google.maps.InfoWindow({ position:thisObj.latlng_} );
		Canvas2.infowindow_.setContent("<div class='infowin-content' style='height:205px'>" + thisObj.msg_ + "</div>");
		Canvas2.infowindow_.open(thisObj.map_);
		Geocoding.getAddress(thisObj.latlng_, "span_");
		*/
			//var_dump(event);
			
		//document.getElementById("Memo").value = 'x=' + event.offsetX + ', y=' + event.offsetY + ', val=' + val + ', r=' + Math.pow(r, 2);	
	}
	
	this.div_.onmouseover = function(event) 
	{

	}

	this.div_.onmouseout = function(event) 
	{
	
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
 

Canvas2.prototype = new google.maps.OverlayView;
Canvas2.infowindow_ = null;
Canvas2.tooltip = null;





// Implement onAdd
Canvas2.prototype.onAdd = function() {
	 
    
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
Canvas2.prototype.onRemove = function() {

	if (this.div_.parentNode !=null)
this.div_.parentNode.removeChild(this.div_);

 // Label is removed from the map, stop updating its position/text.
 for (var i = 0, I = this.listeners_.length; i < I; ++i) {
   google.maps.event.removeListener(this.listeners_[i]);
 }
};

// Implement draw
Canvas2.prototype.draw = function() {
	
	var projection = this.getProjection();
	var pos =   new google.maps.LatLng(this.latlng_.lat(), this.latlng_.lng());
	position = projection.fromLatLngToDivPixel(pos);
	
	this.div_.width = this.div_.width;
	var div = this.div_; 
	var x0 = (position.x - parseFloat(div.width)/2 - 2 );
	var y0 = (position.y - parseFloat(div.height)/2 - 2);
	
	div.style.left = (x0) ;
	div.style.top = (y0) ;
	div.style.display = 'block';
	var ctx = this.ctx_;

	ctx.fillStyle = this.fillcolor_; // ���� ����
	ctx.strokeStyle = this.fillcolor_; // 磥⡴ﮠ
	
	// var img = new Image();
	//img.src = '/i/js_map_img/strelka.png';
	
	ctx.lineWidth = 1.3;
	ctx.save();  
	ctx.moveTo(0, 0);
	var cx = this.cx;
	var cy = this.cy;
	var d = this.d;
	
	/*ctx.translate(cx, cy);
	ctx.rotate(this.angle_* Math.PI / 180);	
	ctx.translate(-cx, -cy);
	ctx.drawImage(img,cx-5, cy-5);
*/
	if (this.type_ == "arrow")
	{
		ctx.translate(cx, cy);
		ctx.rotate(this.angle_* Math.PI / 180);	
		ctx.translate(-cx, -cy);
	
		ctx.moveTo(cx + d , cy + d );
		ctx.lineTo(cx , cy - d+4);
		ctx.lineTo(cx - d , cy + d);
		ctx.moveTo(cx,  cy - d+4);
		ctx.lineTo(cx, cy - d + 20 );
		
		
		
	
	//	ctx.closePath();
		//ctx.fill();
		ctx.stroke();
	}
	else
	{
		ctx.beginPath();
		ctx.arc(cx,cy,3,0,Math.PI*2,true); 
		ctx.closePath();
		ctx.fill();
		ctx.stroke();
	}
	
	
};


	

