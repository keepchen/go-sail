let script = `var pluginsCopy = function(className) {
    var d = document.getElementsByClassName(className),
        quickCopy = function(ele) {
            var d = document.createElement("textarea");
            d.setAttribute("readonly", "readonly");
            d.value = ele.innerText;
            d.style.opacity = "0";
            // d.style.display = "none";
            ele.parentElement.appendChild(d);
            d.select();
            if (d.createTextRange) {
                var c = d.createTextRange();
                c.collapse(true);
                c.moveStart("character", 0);
                c.moveEnd("character", ele.innerText.length);
            } else {
                if (d.setSelectionRange) {
                    d.setSelectionRange(0, ele.innerText.length);
                }
            }
            if (document.execCommand == undefined) {
                d.remove();
                throw new Error("[copy] copy failed, \`execCommand\` not be supported by your browser.");
            }
            var b = document.execCommand("Copy");
            ele.parentElement.removeChild(d);
            if ("getSelection" in window) {
                window.getSelection().removeAllRanges();
            } else {
                document.selection.empty();
            }
            d.remove();
            return b;
        },
        appendNode = function(target) {
            copySpan = document.createElement("span");
            copySpan.style.marginRight = "10px";
            copySpan.innerText = "Copy";
            copySpan.onclick = cp;
            target.parentElement.insertBefore(copySpan, target.parentElement.lastChild);
        },
        cp = function() {
            var target = this.parentElement.firstChild.nextSibling;
            if (quickCopy(target)) {
                this.innerText = "Copied";
                var vm = this;
                setTimeout(function(){
                    vm.innerText = "Copy";
                }, 10000);
            }
        };
    for (let i=0;i<d.length;i++) {
        appendNode(d[i]);
    }
}

setTimeout(function(){ pluginsCopy("http-verb"); }, 5000);
`

if (process.argv.length < 3) {
    console.log("Nothing to do.");
    console.log(`------------------------Usage------------------------
node ${process.argv[1]} path/to/docs.html
`);
    return false;
}

let src = process.argv[2];
let fs = require("fs");

fs.readFile(src, "utf8", function (err,data) {
    if (err) {
        return console.log(err);
    }
    var result = data.replace(/<\/body>/g, "<script>\n"+ script + "<\/script>\n<\/body>");

    fs.writeFile(src, result, "utf8", function (err) {
        if (err) return console.log(err);
    });
});

console.log("Inject redocly copy plugin success!");