const userId = pageData.user_id;
const excelSvg=`<svg class="card-svg-icon" xmlns="http://www.w3.org/2000/svg" viewBox="0 -960 960 960"  fill="#fff"><path d="M120-200v-560q0-33 23.5-56.5T200-840h560q33 0 56.5 23.5T840-760v560q0 33-23.5 56.5T760-120H200q-33 0-56.5-23.5T120-200Zm80-400h560v-160H200v160Zm213 200h134v-120H413v120Zm0 200h134v-120H413v120ZM200-400h133v-120H200v120Zm427 0h133v-120H627v120ZM200-200h133v-120H200v120Zm427 0h133v-120H627v120Z"/></svg>`
const markdownSvg = `<svg class="card-svg-icon" xmlns="http://www.w3.org/2000/svg"  viewBox="0 -960 960 960"  fill="#fff"><path d="m640-360 120-120-42-43-48 48v-125h-60v125l-48-48-42 43 120 120ZM160-160q-33 0-56.5-23.5T80-240v-480q0-33 23.5-56.5T160-800h640q33 0 56.5 23.5T880-720v480q0 33-23.5 56.5T800-160H160Zm0-80h640v-480H160v480Zm0 0v-480 480Zm60-120h60v-180h40v120h60v-120h40v180h60v-200q0-17-11.5-28.5T440-600H260q-17 0-28.5 11.5T220-560v200Z"/></svg>`
const switchSvg = (documentType)=>{
    switch (documentType){
        case 1:
            return markdownSvg;
        case 2:
            return excelSvg;
        default:
            return "";
    }
}
$.ajax({
    url: `/api/user/${userId}/documents`,
    type: 'GET',
    dataType: 'json',
    timeout: 3e4,
}).done(function (data) {
    let documentListHtml = "";
    for (let i = 0; i < data.length; i++) {
        const documentInfo = data[i];
        const documentId = documentInfo.document_id;
        const documentType = documentInfo.document_type;
        const title = documentInfo.title;
        //const author = documentInfo.author;
        const updated = documentInfo.updated;

        documentListHtml += `
<a href="/document/${documentId}" class="document-card">
    <div class="document-card-icon-container">
        <div class="document-card-icon icon-type-${documentType}">${switchSvg(documentType)}</div>
    </div>
    <span class="document-card-title">${title}</span>
    
    <span class="document-card-updated">${updated}</span>
</a>
`;
    }
    $('#document-list').append(documentListHtml);
});
