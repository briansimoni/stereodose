/* 
INSPIRED BY DRIBBLE SHOTS
musics taken from krafta
https://dribbble.com/shots/1478420-Music-app
https://dribbble.com/shots/1249898-Music-Player-Side-Menu
*/


$("#random_track, #repeat_track").click(function(){
  $(this).toggleClass("active-button");
})
$(".play-pause").click(function(){
  if($(this).hasClass("pause-active")){
    $(this).removeClass("pause-active");
    $(this).addClass("play-active");
    $(".fa-play").css("display","inline-block");
    $(".fa-pause").hide();
  }else{
    $(this).addClass("pause-active");
    $(this).removeClass("play-active");
    $(".fa-play").hide();
    $(".fa-pause").show();
  }
});

function secondsToTime(secs){
    secs = Math.round(secs);

    var divisor_for_minutes = secs % (60 * 60);
    var minutes = Math.floor(divisor_for_minutes / 60);

    var divisor_for_seconds = divisor_for_minutes % 60;
    var seconds = Math.ceil(divisor_for_seconds);

    var time = "";
    if(minutes < 10){
      time = "0"+minutes;
    }else{
      time = minutes;
    }
    time += ":";
    if(seconds < 10){
      time += "0"+seconds;
    }else{
      time += seconds;
    }
    return time;
}
/*ADDING MAGIC TO THE PLAYER*/
var track_time = -1;  //in seconds
var total_time = 409; //in seconds
var play = setInterval(function(){updateTime()}, 1000);

function updateTime(){
  track_time++;
  
  //updating timer
  $("#track_time").html(secondsToTime(track_time));
  
  //updating status-progress
  var progress_width = (100*track_time)/total_time;
  $(".status-progress").css("width", progress_width+"%");
  
  //on finish track, go to the next one
  if(track_time >= total_time){
    audio.pause();
    stopPlaying();
    if(playing_track < tracks.track.length){
      playing_track++; 
    }
    //change track here
    playTrack();
  }
}
function stopPlaying() {
    clearInterval(play);
}

var tracks = {"track":[
    {
        "trackName" : "Time",
        "trackArtist" : "Pink Floyd",
        "trackUrl":"https://api.soundcloud.com/tracks/16362406/stream?client_id=5830ff6714c2d3bc7aa3db9947f23231", 
        "trackImage":"https://s3-us-west-2.amazonaws.com/s.cdpn.io/953/dark.jpg", 
        "trackDuration":"409"
    },
    {
        "trackName" : "Paranoid",
        "trackArtist" : "Black Sabbath",
        "trackUrl":"https://api.soundcloud.com/tracks/23256180/stream?client_id=5830ff6714c2d3bc7aa3db9947f23231", 
        "trackImage":"http://upload.wikimedia.org/wikipedia/en/6/64/Black_Sabbath_-_Paranoid.jpg", 
        "trackDuration":"173"
    },
    {
        "trackName" : "I Talk to the Wind",
        "trackArtist" : "King Crimson",
        "trackUrl":"https://api.soundcloud.com/tracks/122200268/stream?client_id=5830ff6714c2d3bc7aa3db9947f23231", 
        "trackImage":"https://consequenceofsound.files.wordpress.com/2013/09/kingcrimson21st.jpg", 
        "trackDuration":"367"
    },
    {
        "trackName" : "Ghost of Perdition",
        "trackArtist" : "Opeth",
        "trackUrl":"https://api.soundcloud.com/tracks/62222301/stream?client_id=5830ff6714c2d3bc7aa3db9947f23231", 
        "trackImage":"http://www.metalsucks.net/wp-content/uploads/2014/11/Opeth-Ghost-Reveries.jpg", 
        "trackDuration":"630"
    },
]};

var audio = new Audio('https://api.soundcloud.com/tracks/16362406/stream?client_id=5830ff6714c2d3bc7aa3db9947f23231');
audio.play();

$(".play-pause").click(function(){
  if($(this).hasClass("play-active")){
    audio.pause();
    stopPlaying();
  }else{
    audio.play();
    var play = setInterval(function(){updateTime()}, 1000);
  }
});

var playing_track = 0; //index of track vector
$(".next-track").click(function(){
  if(playing_track < tracks.track.length){
    playing_track++; 
  }
  audio.pause();
  stopPlaying();
  //change track here
  playTrack();
});
$(".prev-track").click(function(){
  if(playing_track>0){
    playing_track--;
  }
  audio.pause();
  stopPlaying();
  //change track here
  playTrack();
});

function playTrack(){
  track_time = -1;
  total_time = tracks.track[playing_track].trackDuration;
  
  
  //updating images
  $(".background").css("background","url("+tracks.track[playing_track].trackImage+")");
  $(".background").css("background-size","cover");
  $(".disco-cover").css("background","url("+tracks.track[playing_track].trackImage+")");
  $(".disco-cover").css("background-size","cover");
  $(".disco-center").css("background","url("+tracks.track[playing_track].trackImage+")");
  $(".disco-center").css("background-size","cover");
  
  //update track info
  $("#track_name").html(tracks.track[playing_track].trackName);
  $("#track_artist").html(tracks.track[playing_track].trackArtist);
  
  //updating timer
  $("#track_time").html("00:00");
  $("#track_duration").html(secondsToTime(tracks.track[playing_track].trackDuration));
  
  //updating status-progress
  var progress_width = (100*track_time)/total_time;
  $(".status-progress").css("width", progress_width+"%");
  
  //play the audio
  audio = new Audio(tracks.track[playing_track].trackUrl);
  audio.play();
  
  play = setInterval(function(){updateTime()}, 1000);
}