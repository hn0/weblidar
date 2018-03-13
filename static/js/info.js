

function Info(data, stream) {

    this.totalPts  = data.PointCnt || 1;
    this.loadedPts = 1;
    this.progress  = document.createElement( 'canvas' );

    stream.on( 'pointbatch', this.progress_tick.bind( this ) );
    stream.on( 'done', this.end_progress.bind( this ) );

    let infoparts = document.getElementsByClassName( 'info' );

    ['No. pts', 'View center'].forEach(function(lbl, i){
        let span = document.createElement( 'span' );
        span.className = 'lbl';
        span.appendChild( document.createTextNode( lbl ) );
        infoparts[i].appendChild( span );
    });

    // left side total number of items
    let ncnt = document.createElement( 'div' );
    ncnt.appendChild( document.createTextNode( this.totalPts ) );
    infoparts[0].appendChild( ncnt );

    // progress bar stuff, for now keep it simple
    let pwrapper = document.getElementById( 'progress' );
    if( pwrapper ){
        pwrapper.appendChild( this.progress );
    }
    else {
        console.warn( 'where can progress canvas alternately go?' );
    }
};

Info.prototype.start_progress = function()
{
    console.log( 'start drawing progress!' );
};

Info.prototype.end_progress = function()
{
    console.log( 'end drawing progress!' );
};

Info.prototype.progress_tick = function(cnt)
{
    console.log( 'update progress tick', cnt );
};

Info.prototype.set_viewoport = function(viewport)
{

};