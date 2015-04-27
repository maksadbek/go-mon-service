var MonitoringTools =
{
	arr_tools:[],		
	MonitoringInit:function(arr)
	{		
		this.arr_tools = arr;	
		if (arr.length == 0)
		{
			this.arr_tools.push({id:0, radius:5, fillcolor:'#D3005F', addparams:0});
		}
		
	}
		
}
