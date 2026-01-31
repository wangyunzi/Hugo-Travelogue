function loadCssCode(code) {
  var style = document.createElement('style');
  style.type = 'text/css';
  style.rel = 'stylesheet';
  style.appendChild(document.createTextNode(code));
  var head = document.getElementsByTagName('head')[0];
  head.appendChild(style);
}

var allCSS = `
  #bber {
    margin-top: 2em;
    width: 100% !important;
  }
  .bb-timeline {
    padding: 0;
    list-style: none;
  }
  .bb-timeline ul {
    padding: 0;
    margin: 0;
  }
  .bb-timeline ul li {
    position: relative;
    margin: 0 0 1rem 0; /* 减少下方间距 */
    padding-left: 1rem;
    width: 100%;
    list-style: none;
  }
  .bb-timeline ul li .bb-div {
    position: relative;
    padding: 0.5rem 0;
  }
  .bb-timeline ul li .datatime {
    display: flex;
    align-items: center;
    margin-bottom: 0.5rem;
  }
  .bb-timeline ul li .datatime .hy-time {
    font-size: 0.8rem;
    color: #999;
  }
  .bb-timeline ul li .datacont {
    font-size: 1rem;
    line-height: 1.6;
  }
  .bb-timeline ul li .datacont img {
    max-width: 100%;
    border-radius: 0.5rem;
  }
  .bb-timeline ul li .tag-span {
    color: #98c1d9;
  }
  .bb-timeline ul li .datasource a {
    color: #fafafa;
    background: #3b3d42;
    padding: 2px 8px;
    margin: 0 6px 0 0;
    border-radius: 5px;
    font-size: 0.9rem;
    font-weight: 400;
  }
  .bb-load {
    text-align: center;
    margin: 2rem 0;
  }
  .bb-load button {
    padding: 0.5rem 1.5rem;
    font-size: 0.9rem;
    background: none;
    border: 1px solid #475671;
    border-radius: 0.5rem;
    color: #475671;
    cursor: pointer;
    transition: background 0.3s, color 0.3s;
  }
  .bb-load button:hover {
    background: #475671;
    color: #fff;
  }
`;

loadCssCode(allCSS);

var limit = bbMemo.limit;
var memos = bbMemo.memos;
var page = 1, offset = 0, nextLength = 0, nextDom = '';
var bbDom = document.querySelector(bbMemo.domId);
var load = '<div class="bb-load"><button class="load-btn button-load">加载中…</button></div>';
if (bbDom) {
  getFirstList();
  meNums();
  var btn = document.querySelector("button.button-load");
  btn.addEventListener("click", function () {
    btn.textContent = '加载中…';
    updateHTMl(nextDom);
    if (nextLength < limit) {
      document.querySelector("button.button-load").remove();
      return;
    }
    getNextList();
  });
}

function getFirstList() {
  bbDom.insertAdjacentHTML('afterend', load);
  var bbUrl = memos + "api/v1/memo?creatorId=" + bbMemo.creatorId + "&rowStatus=NORMAL&limit=" + limit;
  fetch(bbUrl).then(res => res.json()).then(resdata => {
    updateHTMl(resdata);
    var nowLength = resdata.length;
    if (nowLength < limit) {
      document.querySelector("button.button-load").remove();
      return;
    }
    page++;
    offset = limit * (page - 1);
    getNextList();
  });
}

function getNextList() {
  var bbUrl = memos + "api/v1/memo?creatorId=" + bbMemo.creatorId + "&rowStatus=NORMAL&limit=" + limit + "&offset=" + offset;
  fetch(bbUrl).then(res => res.json()).then(resdata => {
    nextDom = resdata;
    nextLength = nextDom.length;
    page++;
    offset = limit * (page - 1);
    if (nextLength < 1) {
      document.querySelector("button.button-load").remove();
      return;
    }
  });
}

function meNums() {
  var bbLoad = document.querySelector('.bb-load');
  var bbUrl = memos + "api/v1/memo/amount?userId=" + bbMemo.creatorId;
  fetch(bbUrl).then(res => res.json()).then(resdata => {
    if (resdata) {}
  });
}

function updateHTMl(data) {
  var resultAll = "";
  const TAG_REG = /#([^\s#]+?) /g,
    BILIBILI_REG = /<a.*?href="https:\/\/www\.bilibili\.com\/video\/((av[\d]{1,10})|(BV([\w]{10})))\/?".*?>.*<\/a>/g,
    NETEASE_MUSIC_REG = /<a.*?href="https:\/\/music\.163\.com\/.*id=([0-9]+)".*?>.*<\/a>/g,
    QQMUSIC_REG = /<a.*?href="https\:\/\/y\.qq\.com\/.*(\/[0-9a-zA-Z]+)(\.html)?".*?>.*?<\/a>/g,
    QQVIDEO_REG = /<a.*?href="https:\/\/v\.qq\.com\/.*\/([a-z|A-Z|0-9]+)\.html".*?>.*<\/a>/g,
    YOUKU_REG = /<a.*?href="https:\/\/v\.youku\.com\/.*\/id_([a-z|A-Z|0-9|==]+)\.html".*?>.*<\/a>/g,
    YOUTUBE_REG = /<a.*?href="https:\/\/www\.youtube\.com\/watch\?v\=([a-z|A-Z|0-9]{11})\".*?>.*<\/a>/g;
  
  marked.setOptions({
    breaks: true,
    smartypants: true,
    langPrefix: 'language-'
  });

  for (var i = 0; i < data.length; i++) {
    var bbContREG = data[i].content.replace(TAG_REG, "<span class='tag-span'>#$1</span> ");
    bbContREG = marked.parse(bbContREG).replace(BILIBILI_REG, "<div class='video-wrapper'><iframe src='//player.bilibili.com/player.html?bvid=$1&as_wide=1&high_quality=1&danmaku=0' scrolling='no' border='0' frameborder='no' framespacing='0' allowfullscreen='true'></iframe></div>").replace(NETEASE_MUSIC_REG, "<meting-js auto='https://music.163.com/#/song?id=$1'></meting-js>").replace(QQMUSIC_REG, "<meting-js auto='https://y.qq.com/n/yqq/song$1.html'></meting-js>").replace(QQVIDEO_REG, "<div class='video-wrapper'><iframe src='//v.qq.com/iframe/player.html?vid=$1' allowFullScreen='true' frameborder='no'></iframe></div>").replace(YOUKU_REG, "<div class='video-wrapper'><iframe src='https://player.youku.com/embed/$1' frameborder=0 'allowfullscreen'></iframe></div>").replace(YOUTUBE_REG, "<div class='video-wrapper'><iframe src='https://www.youtube.com/embed/$1' frameborder='0' allowfullscreen></iframe></div>");
    
    // 确保时间戳正确
    var createdTs = data[i].createdTs;
    var bbDatetime = new Date(createdTs * 1000).toISOString().replace("T", " ").substring(0, 16);

    var dom = '<li><div class="bb-div"><div class="datatime"><span class="hy-time">' + bbDatetime + '</span></div><div class="datacont">' + bbContREG + '</div><div class="datasource">' + (data[i].resourceList.length ? '<a href="' + data[i].resourceList[0].externalLink + '" target="_blank">' + data[i].resourceList[0].title + '</a>' : '') + '</div></div></li>';
    resultAll += dom;
  }

  bbDom.insertAdjacentHTML('beforeend', '<div class="bb-timeline"><ul>' + resultAll + '</ul></div>');
}
