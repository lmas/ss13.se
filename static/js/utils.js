
function pad(value){
    str = String(value);
    if (str.length == 1) {
        str = "0" + str;
    }
    return str;
}

function format_date(time){
    var ts = time.getFullYear() + "-";
    // +1 because this fucker decided to be zero indexed, unlike the rest.
    ts += pad(time.getUTCMonth()+1) + "-";
    ts += pad(time.getUTCDate()) + " ";
    ts += pad(time.getUTCHours()) + ":";
    ts += pad(time.getUTCMinutes()) + " UTC";
    return ts
}

function show_tooltip(msg, x, y){
    $("#tooltip")
        .html(msg)
        .css({top: y, left: x})
        .fadeIn(200);
}

function hide_tooltip(){
    $("#tooltip").hide();
}

function players_at_tooltip(event, pos, item) {
    if (item) {
        var time = item.datapoint[0];
        var players = item.datapoint[1];
        var timestamp = format_date(new Date(time));
        show_tooltip(players + " players at " + timestamp, item.pageX+5, item.pageY+5);
    } else {
        hide_tooltip();
    }
}

function players_tooltip(event, pos, item) {
    if (item) {
        var players = item.datapoint[1];
        show_tooltip(players + " players", item.pageX, item.pageY);
    } else {
        hide_tooltip();
    }
}

function timestamp_tooltip(time){
    var pos = $("#timestamp").offset();
    show_tooltip(time, pos.left, pos.top+20);
}
