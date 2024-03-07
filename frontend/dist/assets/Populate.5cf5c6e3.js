import{I as l}from"./InstallationStep.916efcd3.js";import{_ as n,m as r,O as h,r as u,o as d,l as c,w as m,e,s as S,v as f}from"./index.2e5144e1.js";const _=[{title:"Select Folder",active:!0},{title:"Populate",active:!1}],v="Specify the path to a folder containing the metadata information to populate the environment with.",g={components:{InstallationStep:l},data(){return{steps:_,tips:v}},computed:{...r({populateState:t=>t.populateState}),navigation(){return{next:{path:"/populate-install/"+this.populateState.id,disabled:this.populateState.path===""||this.populateState.path===null},back:{path:`/environments/${this.populateState.id}`,disabled:!1},cancel:{path:`/environments/${this.populateState.id}`,disabled:!1}}}},methods:{openPathDialog(){h("Select the folder containing the metadata").then(t=>{this.populateState.path=t}).catch(t=>{console.error(t)})}},mounted(){this.$store.commit("resetPopulateState"),this.$route.params.id&&(this.populateState.id=this.$route.params.id,this.populateState.name=this.populateState.id.split("-")[0],this.populateState.version=this.populateState.id.split("-")[1],this.populateState.platform=this.populateState.id.split("-")[2])}},b={class:"populate-main-content-container"},P=e("h1",{class:"populate-title"},"Select Folder",-1),x={class:"populate-main-content"},D={class:"populate-form",autocomplete:"off"},k={class:"populate-form-field"},w=e("label",{for:"name"},"Path:",-1);function I(t,a,$,y,s,o){const i=u("InstallationStep");return d(),c(i,{steps:s.steps,tips:s.tips,navigation:o.navigation},{default:m(()=>[e("div",b,[P,e("div",x,[e("form",D,[e("div",k,[w,S(e("input",{type:"text",id:"path",name:"path","onUpdate:modelValue":a[0]||(a[0]=p=>t.populateState.path=p)},null,512),[[f,t.populateState.path]])])]),e("button",{class:"primary-button",onClick:a[1]||(a[1]=(...p)=>o.openPathDialog&&o.openPathDialog(...p))},"Select folder")])])]),_:1},8,["steps","tips","navigation"])}const F=n(g,[["render",I]]);export{F as default};
