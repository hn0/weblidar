

function Info(data) {

    console.log( 'got the data', data );

    let infoparts = document.getElementsByClassName( 'info' );

    ['No. pts', 'View center'].forEach(function(lbl, i){
        let span = document.createElement( 'span' );
        span.className = 'lbl';
        infoparts[i].appendChild( document.createTextNode( lbl ) );
    });

    // left side total number of items
    let ncnt = document.createElement( 'div' );
    ncnt.appendChild( document.createTextNode( data.PointCnt ) );
    infoparts[0].appendChild( ncnt );

    // right side, a good place to put viewport coordinates into!
    // where to place datasource!


    // progress bar stuff, for now keep it simple
    let pwrapper = document.getElementById( 'progress' );
    let progress = document.createElement( 'canvas' );
    if( pwrapper ){

        let ctx = progress.getContext( '2d' );
        ctx.fillColor = '#f04';
        ctx.fillRect(  0, 0, 15, 15 );

        pwrapper.appendChild( progress );
    }
    else {
        console.warn( 'where can progress canvas alternately go?' );
    }
};

Info.prototype.start_progress = function()
{

};

Info.prototype.progress_tick = function()
{

};

Info.prototype.set_viewoport = function()
{

};