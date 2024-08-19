const userId = pageData.user_id;
$.ajax({
    url: `/api/user/${userId}/documents`,
    type: 'GET',
    dataType: 'json',
    timeout: 3e4,
}).done(function (data) {
    let documentListHtml = "";
    for (let i = 0; i < data.length; i++) {
        const documentItem = data[i];
        documentListHtml += `<br><a href="/document/${documentItem.document_id}">类型：${documentItem.permission_type}，编号：${documentItem.document_id}</a>`;
    }
    $('#document-list').append(documentListHtml);
});
