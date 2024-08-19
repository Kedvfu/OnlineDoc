const markdownText = $("#document-content");
const markdownViewer = $("#document-viewer-container")
const favouriteBtn =$("#document-info-bar-item-favourite")
const saveBtn = $("#document-info-bar-item-save")
const currentTitle = $("#document-info-bar-item-title")
const showMarkdown = () => {
    const markdownHtml = markdown.toHTML(markdownText.val());
    markdownViewer.html(markdownHtml);
    showCodeWindow($("code"));

}

const loadDocument = () => {
    if (pageData.document_id === "new") {
        return;
    }
    $.ajax({
        type: "GET",
        url: `/api/user/${pageData.user_id}/document/${pageData.document_id}/get`,
        dataType: "json",
        timeout: 3e4,
    }).done(function (data) {
        currentTitle.text(data.title);
        markdownText.val(data.content);
        loadDocumentUsers(data)
        showMarkdown();
    })
}


const showCodeWindow = (code) => {
    code.each(function () {
        const codeElement = $(this);
        const lines = codeElement.text().trim().split("\n");
        const language = lines[0];

        const lineNumbers = $.map(lines, (line, index) => `<span class='line-number'>${index+1}</span>`).join("<br>");
        const codeWindow = $(`
<div class='code-window'>
    <div class='code-header'>
        <div class='code-title'>${language}</div>
        <a class='copy-button'>
            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 18 18" fill="none"  stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <rect x="9" y="9" width="8" height="8" rx="2" ry="2"></rect>
            <path d="M5 15H4a2 2 0 0 1-2-2V4a2 2 0 0 1 2-2h9a2 2 0 0 1 2 2v1"></path>
            </svg>
            Copy
        </a>
    </div>
    <div class='code-body'>
        <div class='line-numbers'>${lineNumbers}</div>
        <pre><code>${codeElement.text().substring(language.length+1).trim()}</code></pre>
    </div>
</div>
`);
        codeElement.replaceWith(codeWindow);

    })
}

const addFavourite = () => {

}

const saveDocument =  () =>{
    let postData = {}
    if (currentTitle.text() === pageData.title) {
        postData = {
            "content": markdownText.val()
        }
    }else {
        postData = {
            "title": currentTitle.text(),
            "content": markdownText.val()
        }
    }
    $.ajax({
        type: "POST",
        url: `/api/user/${pageData.user_id}/document/${pageData.document_id}/save`,
        dataType: "application/json",
        data: JSON.stringify(postData),
    })
}

loadDocument();
markdownText.on("input", showMarkdown)
favouriteBtn.on("click", addFavourite)
saveBtn.on("click", saveDocument)
