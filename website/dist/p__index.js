(self["webpackChunk"]=self["webpackChunk"]||[]).push([[866],{46076:function(e,t,r){"use strict";r.r(t),r.d(t,{default:function(){return j}});r(71748);var n=r(90860),a=(r(13062),r(71230)),s=(r(89032),r(15746)),c=(r(58024),r(91894)),i=(r(95300),r(7277)),o=r(8870),d=r(57337),u=r(58638),l=r(67294),h=r(28435),p=r(46925),x=r(85893);function j(e){var t=(0,l.useState)(!1),r=(0,d.Z)(t,2),j=r[0],Z=r[1],f=(0,l.useState)({taskCount:0,runCount:0,schedulerCount:0}),g=(0,d.Z)(f,2),v=g[0],C=g[1];(0,l.useEffect)((()=>{Z(!0),(0,h.HT)("/api/u/dashboard",(function(e,t){Z(!1),200==t&&(C(e.data),k(e.data.chart))}))}),[]);var k=e=>{e=e.map((e=>(200==e.code?e.code="\u6210\u529f":e.code="\u5931\u8d25",(0,o.Z)((0,o.Z)({},e),{},{code:e.code+""}))));var t=new p.kL({container:"container",autoFit:!0,height:500});t.data(e),t.scale({date:{range:[0,1]},count:{nice:!0}}),t.tooltip({showCrosshairs:!0,shared:!0}),t.axis("count",{label:{formatter:e=>e}}),t.line().position("date*count").color("code",(e=>"\u5931\u8d25"==e?"red":"green")).shape("smooth"),t.point().position("date*count").color("code",(e=>"\u5931\u8d25"==e?"red":"green")).shape("circle"),t.render()};return(0,x.jsx)(u._z,{extra:[(0,x.jsx)("a",{href:"/metrics",target:"_blank",children:"Metrics"},"metrics")],footer:[],children:(0,x.jsx)(n.Z,{loading:j,active:!0,paragraph:{rows:18},children:(0,x.jsxs)("div",{className:"site-statistic-demo-card",children:[(0,x.jsxs)(a.Z,{gutter:16,children:[(0,x.jsx)(s.Z,{span:12,children:(0,x.jsx)(c.Z,{children:(0,x.jsx)(i.Z,{title:"\u4efb\u52a1\u6570\u91cf / \u8fd0\u884c\u6570\u91cf",value:v.taskCount+" / "+v.runCount,precision:0})})}),(0,x.jsx)(s.Z,{span:12,children:(0,x.jsx)(c.Z,{children:(0,x.jsx)(i.Z,{title:"\u7d2f\u8ba1\u8c03\u5ea6\u6b21\u6570",value:v.schedulerCount,precision:0})})})]}),(0,x.jsx)(a.Z,{gutter:16,style:{marginTop:40},children:(0,x.jsx)(s.Z,{span:24,children:(0,x.jsx)("div",{id:"container"})})})]})})},"xxz")}}}]);