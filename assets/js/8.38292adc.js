(window.webpackJsonp=window.webpackJsonp||[]).push([[8],{328:function(e,t,o){"use strict";o.r(t);var r=o(33),a=Object(r.a)({},(function(){var e=this,t=e.$createElement,o=e._self._c||t;return o("ContentSlotsDistributor",{attrs:{"slot-key":e.$parent.slotKey}},[o("h1",{attrs:{id:"_1-deploy-ocf-cloud"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#_1-deploy-ocf-cloud"}},[e._v("#")]),e._v(" 1. Deploy OCF Cloud")]),e._v(" "),o("p",[e._v("There are multiple options how to start using / testing the gOCF Cloud. If you're just trying to get in touch with this IoT solution, jump right to the "),o("strong",[e._v("free")]),e._v(" "),o("a",{attrs:{href:"#pluggedin.cloud"}},[e._v("pluggedin.cloud instance")]),e._v(" and onboard your device. In case you want to "),o("strong",[e._v("test")]),e._v(" the system localy and you have the "),o("a",{attrs:{href:"https://docs.docker.com/get-docker/",target:"_blank",rel:"noopener noreferrer"}},[e._v("Docker ready"),o("OutboundLink")],1),e._v(", use our "),o("a",{attrs:{href:"bundle"}},[e._v("Bundle Docker Image")]),e._v(". Last but not least, in case you're already familiar with the gOCF Cloud and you want to deploy production ready system, follow our "),o("a",{attrs:{href:"#kubernetes"}},[e._v("Kubernetes deployment")]),e._v(" using Helm Charts(https://helm.sh/).\nIf you're already familiar with the OCF Cloud and want to deploy full-blown system, go right to the (#some-markdown-heading)")]),e._v(" "),o("h2",{attrs:{id:"try-pluggedin-cloud"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#try-pluggedin-cloud"}},[e._v("#")]),e._v(" Try pluggedin.cloud")]),e._v(" "),o("p",[e._v("Simply visit "),o("a",{attrs:{href:"https://pluggedin.cloud",target:"_blank",rel:"noopener noreferrer"}},[e._v("pluggedin.cloud"),o("OutboundLink")],1),e._v(" and click "),o("code",[e._v("Try")]),e._v(".")]),e._v(" "),o("h2",{attrs:{id:"bundle"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#bundle"}},[e._v("#")]),e._v(" Bundle")]),e._v(" "),o("p",[e._v("Bundle option hosts all gOCF Cloud Services and it's dependencies in a single Docker Image. This solution should be used only for "),o("strong",[e._v("testing purposes")]),e._v(" as the authorization servies is not in full operation.")]),e._v(" "),o("h3",{attrs:{id:"pull-the-image"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#pull-the-image"}},[e._v("#")]),e._v(" Pull the image")]),e._v(" "),o("div",{staticClass:"language-bash extra-class"},[o("pre",{pre:!0,attrs:{class:"language-bash"}},[o("code",[e._v("docker pull ocfcloud/bundle:vnext\ndocker run -d --network"),o("span",{pre:!0,attrs:{class:"token operator"}},[e._v("=")]),e._v("host --name"),o("span",{pre:!0,attrs:{class:"token operator"}},[e._v("=")]),e._v("cloud -t ocfcloud/bundle:vnext\n")])])]),o("h3",{attrs:{id:"remarks"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#remarks"}},[e._v("#")]),e._v(" Remarks")]),e._v(" "),o("ul",[o("li",[e._v("OAuth2.0 Authorization Code is not verified during device onboarding")]),e._v(" "),o("li",[e._v("Cloud2Cloud is not part of the Bundle deployment")])]),e._v(" "),o("h2",{attrs:{id:"kubernetes"}},[o("a",{staticClass:"header-anchor",attrs:{href:"#kubernetes"}},[e._v("#")]),e._v(" Kubernetes")]),e._v(" "),o("p",[e._v("comming soon...")])])}),[],!1,null,null,null);t.default=a.exports}}]);