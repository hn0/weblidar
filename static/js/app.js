
var app = function() {
	this.load_url( 'info' )
	    .then(function(data){
	    	var c = new LidarCanvas( this );
	    })
	    .catch( function(){
	    	console.warn( 'Application unable to start, cannot load the data' );
	    });
};

app.prototype.load_url = function(url, method, data) {

	if( !method ){
		method = "GET";
	}

	var p = new Promise(function(resolve,reject){

		var xhr = new XMLHttpRequest();
		xhr.open( method, url );
		xhr.onload = function(){
			resolve( xhr.responseText );
		};

		xhr.onerror = function() {
			reject();
		}

		xhr.send();
	});

	return p;
};