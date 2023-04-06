package selector

import (
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html"
	"strings"
	"testing"
)

//http://gogodjzhu.com
var htmlStr = `<!DOCTYPE html>
<html lang="zh-CN" prefix="og: http://ogp.me/ns#">
    <head>
        <meta charset="UTF-8">
        <meta name="description" content="GoGoDJZhu - A Self Place">
        <meta name="keywords" content="GoGoDJZhu,DJZhu,djzhu,gogodjzhu">
        <link type="text/css" media="all" href="http://gogodjzhu.com/wp-content/cache/autoptimize/css/autoptimize_d4289a86877ab5f852511ee627fa900a.css" rel="stylesheet" />
        <title>GoGo DJZhu - Everything about life and work.</title>
        <meta name="description" content="Everything about life and work."/>
        <link rel="canonical" href="http://gogodjzhu.com/" />
        <link rel="next" href="http://gogodjzhu.com/index.php/page/2/" />
        <meta property="og:locale" content="zh_CN" />
        <meta property="og:type" content="website" />
        <meta property="og:title" content="GoGo DJZhu - Everything about life and work." />
        <meta property="og:description" content="Everything about life and work." />
        <meta property="og:url" content="http://gogodjzhu.com/" />
        <meta property="og:site_name" content="GoGo DJZhu" />
        <meta name="twitter:card" content="summary_large_image" />
        <meta name="twitter:description" content="Everything about life and work." />
        <meta name="twitter:title" content="GoGo DJZhu - Everything about life and work." />
        <script type='application/ld+json'>{"@context":"https:\/\/schema.org","@type":"WebSite","@id":"#website","url":"http:\/\/gogodjzhu.com\/","name":"GoGo DJZhu","potentialAction":{"@type":"SearchAction","target":"http:\/\/gogodjzhu.com\/?s={search_term_string}","query-input":"required name=search_term_string"}}</script>
        <script type='application/ld+json'>{"@context":"https:\/\/schema.org","@type":"Person","url":"http:\/\/gogodjzhu.com\/","sameAs":[],"@id":"#person","name":"DJ Zhu"}</script>
        <link rel='dns-prefetch' href='//fonts.googleapis.com' />
        <link rel='dns-prefetch' href='//s.w.org' />
        <link rel="alternate" type="application/rss+xml" title="GoGo DJZhu &raquo; Feed" href="http://gogodjzhu.com/index.php/feed/" />
        <link rel="alternate" type="application/rss+xml" title="GoGo DJZhu &raquo; 评论Feed" href="http://gogodjzhu.com/index.php/comments/feed/" />
        <script type="text/javascript">window._wpemojiSettings = {"baseUrl":"https:\/\/s.w.org\/images\/core\/emoji\/11\/72x72\/","ext":".png","svgUrl":"https:\/\/s.w.org\/images\/core\/emoji\/11\/svg\/","svgExt":".svg","source":{"concatemoji":"http:\/\/gogodjzhu.com\/wp-includes\/js\/wp-emoji-release.min.js?ver=4.9.9"}};
			!function(a,b,c){function d(a,b){var c=String.fromCharCode;l.clearRect(0,0,k.width,k.height),l.fillText(c.apply(this,a),0,0);var d=k.toDataURL();l.clearRect(0,0,k.width,k.height),l.fillText(c.apply(this,b),0,0);var e=k.toDataURL();return d===e}function e(a){var b;if(!l||!l.fillText)return!1;switch(l.textBaseline="top",l.font="600 32px Arial",a){case"flag":return!(b=d([55356,56826,55356,56819],[55356,56826,8203,55356,56819]))&&(b=d([55356,57332,56128,56423,56128,56418,56128,56421,56128,56430,56128,56423,56128,56447],[55356,57332,8203,56128,56423,8203,56128,56418,8203,56128,56421,8203,56128,56430,8203,56128,56423,8203,56128,56447]),!b);case"emoji":return b=d([55358,56760,9792,65039],[55358,56760,8203,9792,65039]),!b}return!1}function f(a){var c=b.createElement("script");c.src=a,c.defer=c.type="text/javascript",b.getElementsByTagName("head")[0].appendChild(c)}var g,h,i,j,k=b.createElement("canvas"),l=k.getContext&&k.getContext("2d");for(j=Array("flag","emoji"),c.supports={everything:!0,everythingExceptFlag:!0},i=0;i
            <j.length;i++)c.supports[j[i]]=e(j[i]),c.supports.everything=c.supports.everything&&c.supports[j[i]],"flag"!==j[i]&&(c.supports.everythingExceptFlag=c.supports.everythingExceptFlag&&c.supports[j[i]]);c.supports.everythingExceptFlag=c.supports.everythingExceptFlag&&!c.supports.flag,c.DOMReady=!1,c.readyCallback=function(){c.DOMReady=!0},c.supports.everything||(h=function(){c.readyCallback()},b.addEventListener?(b.addEventListener("DOMContentLoaded",h,!1),a.addEventListener("load",h,!1)):(a.attachEvent("onload",h),b.attachEvent("onreadystatechange",function(){"complete"===b.readyState&&c.readyCallback()})),g=c.source||{},g.concatemoji?f(g.concatemoji):g.wpemoji&&g.twemoji&&(f(g.twemoji),f(g.wpemoji)))}(window,document,window._wpemojiSettings);
            </script>
            <link rel='stylesheet' id='ogbbblog-oswald-css'  href='https://fonts.googleapis.com/css?family=Oswald%3A300&#038;ver=4.9.9' type='text/css' media='all' />
            <link rel='https://api.w.org/' href='http://gogodjzhu.com/index.php/wp-json/' />
            <link rel="EditURI" type="application/rsd+xml" title="RSD" href="http://gogodjzhu.com/xmlrpc.php?rsd" />
            <link rel="wlwmanifest" type="application/wlwmanifest+xml" href="http://gogodjzhu.com/wp-includes/wlwmanifest.xml" />
            <meta name="generator" content="WordPress 4.9.9" />
            <link id="favicon" href="http://gogodjzhu.com/favicon.ico" rel="icon" type="image/x-icon" />
        </head>
        <body class="home blog custom-background" >
            <div id="wrapper">
                <div id="header">
                    <div id="site">
                        <div id="sitename">
                            <a href="http://gogodjzhu.com/"> GoGo DJZhu</a>
                        </div>
                        <div id="slogan">Everything about life and work.</div>
                    </div>
                    <div id="banner"></div>
                </div>
                <div id="menu">
                    <div class="menu-top-menu-container">
                        <ul id="menu-top-menu" class="menu">
                            <li id="menu-item-15" class="menu-item menu-item-type-custom menu-item-object-custom current-menu-item current_page_item menu-item-home menu-item-15">
                                <a href="http://gogodjzhu.com">Home</a>
                            </li>
                            <li id="menu-item-178" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-has-children menu-item-178">
                                <a href="http://gogodjzhu.com/index.php/category/music/">Music</a>
                                <ul class="sub-menu">
                                    <li id="menu-item-229" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-229">
                                        <a href="http://gogodjzhu.com/index.php/category/music/tutorial/">Tutorial</a>
                                    </li>
                                    <li id="menu-item-227" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-227">
                                        <a href="http://gogodjzhu.com/index.php/category/music/lick/">Lick</a>
                                    </li>
                                    <li id="menu-item-228" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-228">
                                        <a href="http://gogodjzhu.com/index.php/category/music/resources/">Resource</a>
                                    </li>
                                </ul>
                            </li>
                            <li id="menu-item-179" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-has-children menu-item-179">
                                <a href="http://gogodjzhu.com/index.php/category/code/">Code</a>
                                <ul class="sub-menu">
                                    <li id="menu-item-237" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-237">
                                        <a href="http://gogodjzhu.com/index.php/category/code/machine-learning/">Machine Learning</a>
                                    </li>
                                </ul>
                            </li>
                            <li id="menu-item-180" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-180">
                                <a href="http://gogodjzhu.com/index.php/category/read/">Book</a>
                            </li>
                            <li id="menu-item-177" class="menu-item menu-item-type-taxonomy menu-item-object-category menu-item-177">
                                <a href="http://gogodjzhu.com/index.php/category/something/">Blablabla</a>
                            </li>
                        </ul>
                    </div>
                </div>
                <div id="main">
                    <div id="content">
                        <div class="post">
                            <div class="post-img">
                                <a href="http://gogodjzhu.com/index.php/read/371/"></a>
                            </div>
                            <div class="post-c">
                                <div class="post-title">
                                    <a href="http://gogodjzhu.com/index.php/read/371/">《送你一颗子弹》——读后随笔</a>
                                </div>
                                <div class="post-footer">
                                    <span class="post-date">2018-09-20</span>
                                    <span class="post-author">DJZhu</span>
                                    <span class="post-comment">Leave a comment</span>
                                    <span class="post-category">
                                        <a href="http://gogodjzhu.com/index.php/category/read/" rel="category tag">Book</a>
                                    </span>
                                </div>
                                <div class="post-excerpt">
                                    <p>《送你一颗子弹》记述的是一个女博士的生活，在普通人看起来再简单不过的一 [&#8230;] 
                                        <a href="http://gogodjzhu.com/index.php/read/371/" class="ReadMore">Read More</a>
                                    </p>
                                </div>
                            </div>
                        </div>
                        <div class="post">
                            <div class="post-img">
                                <a href="http://gogodjzhu.com/index.php/something/365/"></a>
                            </div>
                            <div class="post-c">
                                <div class="post-title">
                                    <a href="http://gogodjzhu.com/index.php/something/365/">西班牙内战——当世界还年轻的时候</a>
                                </div>
                                <div class="post-footer">
                                    <span class="post-date">2018-09-12</span>
                                    <span class="post-author">DJZhu</span>
                                    <span class="post-comment">Leave a comment</span>
                                    <span class="post-category">
                                        <a href="http://gogodjzhu.com/index.php/category/something/" rel="category tag">Blablabla</a>
                                    </span>
                                </div>
                                <div class="post-excerpt">
                                    <p>曾经在一篇文章上看到一种说法，历史走到今天，已达尽头，民主法治的西方式 [&#8230;] 
                                        <a href="http://gogodjzhu.com/index.php/something/365/" class="ReadMore">Read More</a>
                                    </p>
                                </div>
                            </div>
                        </div>
                        <div class="post">
                            <div class="post-img">
                                <a href="http://gogodjzhu.com/index.php/something/353/"></a>
                            </div>
                            <div class="post-c">
                                <div class="post-title">
                                    <a href="http://gogodjzhu.com/index.php/something/353/">Practices in memorising the scales and the chords</a>
                                </div>
                                <div class="post-footer">
                                    <span class="post-date">2018-07-09</span>
                                    <span class="post-author">ATM</span>
                                    <span class="post-comment">Leave a comment</span>
                                    <span class="post-category">
                                        <a href="http://gogodjzhu.com/index.php/category/something/" rel="category tag">Blablabla</a>
                                    </span>
                                </div>
                                <div class="post-excerpt">
                                    <p>Here are some simple tips in how to [&#8230;] 
                                        <a href="http://gogodjzhu.com/index.php/something/353/" class="ReadMore">Read More</a>
                                    </p>
                                </div>
                            </div>
                        </div>
                        <div class="post">
                            <div class="post-img">
                                <a href="http://gogodjzhu.com/index.php/code/tools/347/"></a>
                            </div>
                            <div class="post-c">
                                <div class="post-title">
                                    <a href="http://gogodjzhu.com/index.php/code/tools/347/">命令行读取hdfs集群的全部配置</a>
                                </div>
                                <div class="post-footer">
                                    <span class="post-date">2018-07-05</span>
                                    <span class="post-author">DJZhu</span>
                                    <span class="post-comment">Leave a comment</span>
                                    <span class="post-category">
                                        <a href="http://gogodjzhu.com/index.php/category/code/tools/" rel="category tag">Tools</a>
                                    </span>
                                </div>
                                <div class="post-excerpt">
                                    <p>编译以下代码并打成jar包: import org.apache.ha [&#8230;] 
                                        <a href="http://gogodjzhu.com/index.php/code/tools/347/" class="ReadMore">Read More</a>
                                    </p>
                                </div>
                            </div>
                        </div>
                        <div class="post">
                            <div class="post-img">
                                <a href="http://gogodjzhu.com/index.php/code/basic/329/"></a>
                            </div>
                            <div class="post-c">
                                <div class="post-title">
                                    <a href="http://gogodjzhu.com/index.php/code/basic/329/">使用Instrumentation计算java对象大小</a>
                                </div>
                                <div class="post-footer">
                                    <span class="post-date">2018-07-03</span>
                                    <span class="post-author">DJZhu</span>
                                    <span class="post-comment">Leave a comment</span>
                                    <span class="post-category">
                                        <a href="http://gogodjzhu.com/index.php/category/code/basic/" rel="category tag">Basic</a>
                                    </span>
                                </div>
                                <div class="post-excerpt">
                                    <p>对象大小的计算 正如笔者看的这篇文章所描述的一样，当我们试图获取一个J [&#8230;] 
                                        <a href="http://gogodjzhu.com/index.php/code/basic/329/" class="ReadMore">Read More</a>
                                    </p>
                                </div>
                            </div>
                        </div>
                        <div class="post">
                            <div class="post-img">
                                <a href="http://gogodjzhu.com/index.php/code/tools/326/"></a>
                            </div>
                            <div class="post-c">
                                <div class="post-title">
                                    <a href="http://gogodjzhu.com/index.php/code/tools/326/">Astah Professional 7.2.0/1ff236 破解工具</a>
                                </div>
                                <div class="post-footer">
                                    <span class="post-date">2018-07-01</span>
                                    <span class="post-author">DJZhu</span>
                                    <span class="post-comment">Leave a comment</span>
                                    <span class="post-category">
                                        <a href="http://gogodjzhu.com/index.php/category/code/tools/" rel="category tag">Tools</a>
                                    </span>
                                </div>
                                <div class="post-excerpt">
                                    <p>Astah Professional 7.2.0/1ff236 破解工 [&#8230;] 
                                        <a href="http://gogodjzhu.com/index.php/code/tools/326/" class="ReadMore">Read More</a>
                                    </p>
                                </div>
                            </div>
                        </div>
                        <div class="post">
                            <div class="post-img">
                                <a href="http://gogodjzhu.com/index.php/music/resources/322/"></a>
                            </div>
                            <div class="post-c">
                                <div class="post-title">
                                    <a href="http://gogodjzhu.com/index.php/music/resources/322/">Ours Samplus &#8211; Deep Inside</a>
                                </div>
                                <div class="post-footer">
                                    <span class="post-date">2018-06-11</span>
                                    <span class="post-author">DJZhu</span>
                                    <span class="post-comment">Leave a comment</span>
                                    <span class="post-category">
                                        <a href="http://gogodjzhu.com/index.php/category/music/resources/" rel="category tag">Resource</a>
                                    </span>
                                </div>
                                <div class="post-excerpt"></div>
                            </div>
                        </div>
                        <div class="post">
                            <div class="post-img">
                                <a href="http://gogodjzhu.com/index.php/music/resources/316/"></a>
                            </div>
                            <div class="post-c">
                                <div class="post-title">
                                    <a href="http://gogodjzhu.com/index.php/music/resources/316/">All of me &#8211; [The jazz real book series]</a>
                                </div>
                                <div class="post-footer">
                                    <span class="post-date">2018-06-05</span>
                                    <span class="post-author">DJZhu</span>
                                    <span class="post-comment">Leave a comment</span>
                                    <span class="post-category">
                                        <a href="http://gogodjzhu.com/index.php/category/music/resources/" rel="category tag">Resource</a>
                                    </span>
                                </div>
                                <div class="post-excerpt"></div>
                            </div>
                        </div>
                        <div class="post">
                            <div class="post-img">
                                <a href="http://gogodjzhu.com/index.php/something/305/"></a>
                            </div>
                            <div class="post-c">
                                <div class="post-title">
                                    <a href="http://gogodjzhu.com/index.php/something/305/"># 嘘！老猪教你如何冲浪~</a>
                                </div>
                                <div class="post-footer">
                                    <span class="post-date">2018-05-17</span>
                                    <span class="post-author">DJZhu</span>
                                    <span class="post-comment">Leave a comment</span>
                                    <span class="post-category">
                                        <a href="http://gogodjzhu.com/index.php/category/something/" rel="category tag">Blablabla</a>
                                    </span>
                                </div>
                                <div class="post-excerpt">
                                    <p>下载安装对应的shadowsock客户端: Shadowsocks-a [&#8230;] 
                                        <a href="http://gogodjzhu.com/index.php/something/305/" class="ReadMore">Read More</a>
                                    </p>
                                </div>
                            </div>
                        </div>
                        <div class="post">
                            <div class="post-img">
                                <a href="http://gogodjzhu.com/index.php/code/303/"></a>
                            </div>
                            <div class="post-c">
                                <div class="post-title">
                                    <a href="http://gogodjzhu.com/index.php/code/303/">很酷的一些东西</a>
                                </div>
                                <div class="post-footer">
                                    <span class="post-date">2018-05-16</span>
                                    <span class="post-author">DJZhu</span>
                                    <span class="post-comment">Leave a comment</span>
                                    <span class="post-category">
                                        <a href="http://gogodjzhu.com/index.php/category/code/" rel="category tag">Code</a>
                                    </span>
                                </div>
                                <div class="post-excerpt">
                                    <p>命令行 curl http://wttr.in/GuangZhou 一 [&#8230;] 
                                        <a href="http://gogodjzhu.com/index.php/code/303/" class="ReadMore">Read More</a>
                                    </p>
                                </div>
                            </div>
                        </div>
                        <nav class="navigation pagination" role="navigation">
                            <h2 class="screen-reader-text">文章导航</h2>
                            <div class="nav-links">
                                <span aria-current='page' class='page-numbers current'>1</span>
                                <a class='page-numbers' href='http://gogodjzhu.com/index.php/page/2/'>2</a>
                                <a class="next page-numbers" href="http://gogodjzhu.com/index.php/page/2/">Next</a>
                            </div>
                        </nav>
                    </div>
                    <div id="sidebar">
                        <div class="widget">
                            <div id="calendar_wrap" class="calendar_wrap">
                                <table id="wp-calendar">
                                    <caption>2018年十二月</caption>
                                    <thead>
                                        <tr>
                                            <th scope="col" title="星期一">一</th>
                                            <th scope="col" title="星期二">二</th>
                                            <th scope="col" title="星期三">三</th>
                                            <th scope="col" title="星期四">四</th>
                                            <th scope="col" title="星期五">五</th>
                                            <th scope="col" title="星期六">六</th>
                                            <th scope="col" title="星期日">日</th>
                                        </tr>
                                    </thead>
                                    <tfoot>
                                        <tr>
                                            <td colspan="3" id="prev">
                                                <a href="http://gogodjzhu.com/index.php/date/2018/09/">&laquo; 9月</a>
                                            </td>
                                            <td class="pad">&nbsp;</td>
                                            <td colspan="3" id="next" class="pad">&nbsp;</td>
                                        </tr>
                                    </tfoot>
                                    <tbody>
                                        <tr>
                                            <td colspan="5" class="pad">&nbsp;</td>
                                            <td>1</td>
                                            <td>2</td>
                                        </tr>
                                        <tr>
                                            <td>3</td>
                                            <td>4</td>
                                            <td>5</td>
                                            <td>6</td>
                                            <td>7</td>
                                            <td>8</td>
                                            <td>9</td>
                                        </tr>
                                        <tr>
                                            <td>10</td>
                                            <td>11</td>
                                            <td>12</td>
                                            <td>13</td>
                                            <td>14</td>
                                            <td>15</td>
                                            <td>16</td>
                                        </tr>
                                        <tr>
                                            <td>17</td>
                                            <td>18</td>
                                            <td>19</td>
                                            <td>20</td>
                                            <td id="today">21</td>
                                            <td>22</td>
                                            <td>23</td>
                                        </tr>
                                        <tr>
                                            <td>24</td>
                                            <td>25</td>
                                            <td>26</td>
                                            <td>27</td>
                                            <td>28</td>
                                            <td>29</td>
                                            <td>30</td>
                                        </tr>
                                        <tr>
                                            <td>31</td>
                                            <td class="pad" colspan="6">&nbsp;</td>
                                        </tr>
                                    </tbody>
                                </table>
                            </div>
                        </div>
                        <div class="widget">
                            <h2 class="widget-title">文章归档</h2>
                            <ul>
                                <li>
                                    <a href='http://gogodjzhu.com/index.php/date/2018/09/'>2018年九月</a>&nbsp;(2)
                                </li>
                                <li>
                                    <a href='http://gogodjzhu.com/index.php/date/2018/07/'>2018年七月</a>&nbsp;(4)
                                </li>
                                <li>
                                    <a href='http://gogodjzhu.com/index.php/date/2018/06/'>2018年六月</a>&nbsp;(2)
                                </li>
                                <li>
                                    <a href='http://gogodjzhu.com/index.php/date/2018/05/'>2018年五月</a>&nbsp;(3)
                                </li>
                                <li>
                                    <a href='http://gogodjzhu.com/index.php/date/2018/04/'>2018年四月</a>&nbsp;(7)
                                </li>
                            </ul>
                        </div>
                        <div class="widget">
                            <h2 class="widget-title">Tags</h2>
                            <div class="tagcloud">
                                <a href="http://gogodjzhu.com/index.php/tag/blues/" class="tag-cloud-link tag-link-36 tag-link-position-1" style="font-size: 8pt;" aria-label="blues (1个项目)">blues
                                    <span class="tag-link-count"> (1)</span>
                                </a>
                                <a href="http://gogodjzhu.com/index.php/tag/golang/" class="tag-cloud-link tag-link-27 tag-link-position-2" style="font-size: 8pt;" aria-label="golang (1个项目)">golang
                                    <span class="tag-link-count"> (1)</span>
                                </a>
                                <a href="http://gogodjzhu.com/index.php/tag/guitar/" class="tag-cloud-link tag-link-29 tag-link-position-3" style="font-size: 8pt;" aria-label="guitar (1个项目)">guitar
                                    <span class="tag-link-count"> (1)</span>
                                </a>
                                <a href="http://gogodjzhu.com/index.php/tag/hadoop/" class="tag-cloud-link tag-link-43 tag-link-position-4" style="font-size: 8pt;" aria-label="hadoop (1个项目)">hadoop
                                    <span class="tag-link-count"> (1)</span>
                                </a>
                                <a href="http://gogodjzhu.com/index.php/tag/jazz/" class="tag-cloud-link tag-link-39 tag-link-position-5" style="font-size: 8pt;" aria-label="jazz (1个项目)">jazz
                                    <span class="tag-link-count"> (1)</span>
                                </a>
                                <a href="http://gogodjzhu.com/index.php/tag/jazz-hiphop/" class="tag-cloud-link tag-link-41 tag-link-position-6" style="font-size: 8pt;" aria-label="jazz hiphop (1个项目)">jazz hiphop
                                    <span class="tag-link-count"> (1)</span>
                                </a>
                                <a href="http://gogodjzhu.com/index.php/tag/licks/" class="tag-cloud-link tag-link-28 tag-link-position-7" style="font-size: 22pt;" aria-label="licks (2个项目)">licks
                                    <span class="tag-link-count"> (2)</span>
                                </a>
                                <a href="http://gogodjzhu.com/index.php/tag/music/" class="tag-cloud-link tag-link-30 tag-link-position-8" style="font-size: 8pt;" aria-label="music (1个项目)">music
                                    <span class="tag-link-count"> (1)</span>
                                </a>
                                <a href="http://gogodjzhu.com/index.php/tag/tcp-udp/" class="tag-cloud-link tag-link-38 tag-link-position-9" style="font-size: 8pt;" aria-label="TCP/UDP (1个项目)">TCP/UDP
                                    <span class="tag-link-count"> (1)</span>
                                </a>
                                <a href="http://gogodjzhu.com/index.php/tag/video/" class="tag-cloud-link tag-link-40 tag-link-position-10" style="font-size: 8pt;" aria-label="video (1个项目)">video
                                    <span class="tag-link-count"> (1)</span>
                                </a>
                                <a href="http://gogodjzhu.com/index.php/tag/%e5%a5%87%e6%8a%80%e6%b7%ab%e5%b7%a7/" class="tag-cloud-link tag-link-44 tag-link-position-11" style="font-size: 8pt;" aria-label="奇技淫巧 (1个项目)">奇技淫巧
                                    <span class="tag-link-count"> (1)</span>
                                </a>
                            </div>
                        </div>
                    </div>
                </div>
                <div id="footer">
                    <div id="cizgi"></div>
                </div>
            </div>
            <script type="text/javascript" defer src="http://gogodjzhu.com/wp-content/cache/autoptimize/js/autoptimize_a2dea89aa2e9c4e45d3c5e6a92c2ff4a.js"></script>
        </body>
    </html>`

func getNode() *html.Node {
	reader := strings.NewReader(htmlStr)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		panic(err)
	}
	return doc.Get(0)
}

func getNodes() []*html.Node {
	reader := strings.NewReader(htmlStr)
	doc, err := goquery.NewDocumentFromReader(reader)
	if err != nil {
		panic(err)
	}
	return doc.Nodes
}

func TestLinkSelector_Select(t *testing.T) {
	expected := "Select() can not apply to LinkSelector"
	defer func() {
		if recover().(string) != expected {
			t.Error("failed test @ TestLinkSelector_Select, should not be implemented")
		}
	}()
	linkSelector := LinkSelector{}
	linkSelector.Select(getNode())
}

func TestLinkSelector_SelectList(t *testing.T) {
	expected := 67
	linkSelector := LinkSelector{}
	actual := len(linkSelector.SelectList(getNode()))
	if actual != expected {
		t.Errorf("failed test @ TestLinkSelector_SelectList, expecteds:%d, actuals:%d", expected, actual)
	}
}

func TestLinkSelector_SelectNode(t *testing.T) {
	expected := "SelectNode() can not apply to LinkSelector"
	defer func() {
		if recover().(string) != expected {
			t.Error("failed test @ TestLinkSelector_Select, should not be implemented")
		}
	}()
	linkSelector := LinkSelector{}
	linkSelector.SelectNode(getNode())
}

func TestLinkSelector_SelectNodeList(t *testing.T) {
	expected := "SelectNode() can not apply to LinkSelector"
	defer func() {
		if recover().(string) != expected {
			t.Error("failed test @ TestLinkSelector_Select, should not be implemented")
		}
	}()
	linkSelector := LinkSelector{}
	linkSelector.SelectNode(getNode())
}

func TestCssSelector_Select(t *testing.T) {
	//# specified id SelectorText
	expected := "Everything about life and work."
	cssSelector := CssSelector{SelectorText: "div#slogan"}
	actual := cssSelector.Select(getNode())
	if expected != actual {
		t.Errorf("failed test @ TestCssSelector_Select, expecteds:%s, actuals:%s", expected, actual)
	}

	//# specified css SelectorText
	expected = `<a href="http://gogodjzhu.com/index.php/read/371/">《送你一颗子弹》——读后随笔</a>`
	cssSelector = CssSelector{SelectorText: ".post-title"}
	actual = cssSelector.Select(getNode())
	if expected != actual {
		t.Errorf("failed test @ TestCssSelector_Select, expecteds:%s, actuals:%s", expected, actual)
	}

	//# specified css SelectorText and Attribute
	expected = "http://gogodjzhu.com/index.php/read/371/"
	cssSelector = CssSelector{SelectorText: ".post-title > a", AttrName: "href"}
	actual = cssSelector.Select(getNode())
	if expected != actual {
		t.Errorf("failed test @ TestCssSelector_Select, expecteds:%s, actuals:%s", expected, actual)
	}
}

func TestCssSelector_SelectList(t *testing.T) {
	//# specified css selectorText
	expectedStr :=
		`<a href="http://gogodjzhu.com/index.php/read/371/">《送你一颗子弹》——读后随笔</a>
<a href="http://gogodjzhu.com/index.php/something/365/">西班牙内战——当世界还年轻的时候</a>
<a href="http://gogodjzhu.com/index.php/something/353/">Practices in memorising the scales and the chords</a>
<a href="http://gogodjzhu.com/index.php/code/tools/347/">命令行读取hdfs集群的全部配置</a>
<a href="http://gogodjzhu.com/index.php/code/basic/329/">使用Instrumentation计算java对象大小</a>
<a href="http://gogodjzhu.com/index.php/code/tools/326/">Astah Professional 7.2.0/1ff236 破解工具</a>
<a href="http://gogodjzhu.com/index.php/music/resources/322/">Ours Samplus – Deep Inside</a>
<a href="http://gogodjzhu.com/index.php/music/resources/316/">All of me – [The jazz real book series]</a>
<a href="http://gogodjzhu.com/index.php/something/305/"># 嘘！老猪教你如何冲浪~</a>
<a href="http://gogodjzhu.com/index.php/code/303/">很酷的一些东西</a>`
	expectedArr := strings.Split(expectedStr, "\n")
	cssSelector := CssSelector{SelectorText: ".post-title"}
	actualArr := cssSelector.SelectList(getNode())
	for i, actual := range actualArr {
		if expectedArr[i] != actual {
			t.Errorf("failed test @ TestCssSelector_SelectList, expecteds:%s, actuals:%s", expectedArr[i], actual)
		}
	}

	//# specified css selectorText and AttrName
	expectedStr =
		`http://gogodjzhu.com/index.php/read/371/
http://gogodjzhu.com/index.php/something/365/
http://gogodjzhu.com/index.php/something/353/
http://gogodjzhu.com/index.php/code/tools/347/
http://gogodjzhu.com/index.php/code/basic/329/
http://gogodjzhu.com/index.php/code/tools/326/
http://gogodjzhu.com/index.php/music/resources/322/
http://gogodjzhu.com/index.php/music/resources/316/
http://gogodjzhu.com/index.php/something/305/
http://gogodjzhu.com/index.php/code/303/`
	expectedArr = strings.Split(expectedStr, "\n")
	cssSelector = CssSelector{SelectorText: ".post-title > a", AttrName: "href"}
	actualArr = cssSelector.SelectList(getNode())
	for i, actual := range actualArr {
		if expectedArr[i] != actual {
			t.Errorf("failed test @ TestCssSelector_SelectList, expecteds:%s, actuals:%s", expectedArr[i], actual)
		}
	}
}

func TestCssSelector_SelectNode(t *testing.T) {
	//# specified id SelectorText
	expected := "Everything about life and work."
	cssSelector := CssSelector{SelectorText: "div#slogan"}
	actual := goquery.NewDocumentFromNode(cssSelector.SelectNode(getNode())).Text()
	if expected != actual {
		t.Errorf("failed test @ TestCssSelector_SelectNode, expecteds:%s, actuals:%s", expected, actual)
	}

	//# specified css SelectorText
	expected = `<a href="http://gogodjzhu.com/index.php/read/371/">《送你一颗子弹》——读后随笔</a>`
	cssSelector = CssSelector{SelectorText: ".post-title"}
	htmlStr, err := goquery.NewDocumentFromNode(cssSelector.SelectNode(getNode())).Html()
	if err != nil {
		panic(err)
	}
	actual = strings.TrimSpace(htmlStr)
	if expected != actual {
		t.Errorf("failed test @ TestCssSelector_SelectNode, expecteds:%s, actuals:%s", expected, actual)
	}

	//# specified attribute
	//NOT SUPPORT
}

func TestCssSelector_SelectNodeList(t *testing.T) {
	//# specified css selectorText
	expectedStr :=
		`<a href="http://gogodjzhu.com/index.php/read/371/">《送你一颗子弹》——读后随笔</a>
<a href="http://gogodjzhu.com/index.php/something/365/">西班牙内战——当世界还年轻的时候</a>
<a href="http://gogodjzhu.com/index.php/something/353/">Practices in memorising the scales and the chords</a>
<a href="http://gogodjzhu.com/index.php/code/tools/347/">命令行读取hdfs集群的全部配置</a>
<a href="http://gogodjzhu.com/index.php/code/basic/329/">使用Instrumentation计算java对象大小</a>
<a href="http://gogodjzhu.com/index.php/code/tools/326/">Astah Professional 7.2.0/1ff236 破解工具</a>
<a href="http://gogodjzhu.com/index.php/music/resources/322/">Ours Samplus – Deep Inside</a>
<a href="http://gogodjzhu.com/index.php/music/resources/316/">All of me – [The jazz real book series]</a>
<a href="http://gogodjzhu.com/index.php/something/305/"># 嘘！老猪教你如何冲浪~</a>
<a href="http://gogodjzhu.com/index.php/code/303/">很酷的一些东西</a>`
	expectedArr := strings.Split(expectedStr, "\n")
	cssSelector := CssSelector{SelectorText: ".post-title"}
	actualArr := cssSelector.SelectNodeList(getNode())
	for i, actual := range actualArr {
		htmlStr, err := goquery.NewDocumentFromNode(actual).Html()
		if err != nil {
			panic(err)
		}
		actualStr := strings.TrimSpace(htmlStr)
		if expectedArr[i] != actualStr {
			t.Errorf("failed test @ TestCssSelector_SelectNodeList, expecteds:%s, actuals:%s", expectedArr[i], actualStr)
		}
	}
}

func TestRegexSelector_SelectString(t *testing.T) {
	expected := "abcabc"
	regexSelector, _ := NewRegexSelector("(abc)+")
	actual := regexSelector.SelectString("abcabc")
	if expected != actual {
		t.Errorf("failed test @ TestRegexSelector_SelectString, expecteds:%s, actuals:%s", expected, actual)
	}
}

func TestRegexSelector_SelectStringList(t *testing.T) {
	expected := []string{"abc", "abc"}
	regexSelector, _ := NewRegexSelector("(abc)")
	actualArr := regexSelector.SelectStringList("abcabc")
	for i, actual := range actualArr {
		if expected[i] != actual {
			t.Errorf("failed test @ TestRegexSelector_SelectStringList, expecteds:%s, actuals:%s", expected, actual)
		}
	}
}

func TestReplaceSelector_SelectString(t *testing.T) {
	expected := "abcabc"
	regexSelector, _ := NewReplaceSelector("(123)", "abc")
	actual := regexSelector.SelectString("abc123")
	if expected != actual {
		t.Errorf("failed test @ TestReplaceSelector_SelectString, expecteds:%s, actuals:%s", expected, actual)
	}
}

func TestReplaceSelector_SelectStringList(t *testing.T) {
	expected := "SelectStringList() can not apply to ReplaceSelector"
	defer func() {
		if recover().(string) != expected {
			t.Error("failed test @ TestReplaceSelector_SelectStringList, should not be implemented")
		}
	}()
	regexSelector, _ := NewReplaceSelector("(123)", "abc")
	regexSelector.SelectStringList("abc123")
}
