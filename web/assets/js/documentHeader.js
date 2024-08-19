const documentHeaderMenu = $('#document-header-menu');
const documentHeaderMenuBtn = $('#document-info-bar-item-menu');
const documentHeaderMenuShareLinkCreate = $("#document-header-menu-item-share-link-create")
const documentHeaderMenuShareLinkTextArea = $("#document-header-menu-item-share-link")
const documentHeaderMenuShareLinkCopy = $("#document-header-menu-item-share-link-copy")
let menuOpen = false;
const changeMenuState = () => {
    if (menuOpen) {
        documentHeaderMenu.removeClass('menu-open');
        documentHeaderMenu.addClass('menu-closed')
        menuOpen = false;
    } else {
        documentHeaderMenu.addClass('menu-open');
        documentHeaderMenu.removeClass('menu-closed')
        menuOpen = true;
    }
}
const getDocumentShareLink = () => {
    if (pageData.permissionType) {
        $.ajax({
            url: `/api/user/${pageData.userId}/document/${pageData.documentId}/link`,
            type: 'GET',
            timeout: 3e4,
        }).done((data) => {
            const link = data.link;
            documentHeaderMenuShareLinkTextArea.val("http://127.0.0.1:8080/share/"+link);
        });
    }
}
const copyDocumentShareLink = () => {
    documentHeaderMenuShareLinkTextArea.select();
    document.execCommand("copy");

}
const getAllDocumentUsers = () => {
    $.ajax({

    })
}
const loadDocumentUsers = (data) => {
    const users = data.documentUsers;
    for (let i = 0; i < users.length; i++) {
        const user = users[i];
        const userId = user.user_id;
        const permissionType = user.permission_type;
        let permissionTypeContent = "";
        if (permissionType){
            permissionTypeContent ="可编辑"
        } else {
            permissionTypeContent ="只读"
        }
        const permissionItem = $(`<ul>${userId}: (${permissionTypeContent})</ul>`)
        $("#document-header-menu-item-share-link-users").append(permissionItem)
    }
}
documentHeaderMenuShareLinkCreate.on('click',getDocumentShareLink)
documentHeaderMenuBtn.on('click', () => {
    changeMenuState();
});
documentHeaderMenuShareLinkCopy.on('click',copyDocumentShareLink)