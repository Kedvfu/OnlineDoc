const documentHeaderMenu = $('#document-header-menu');
const documentHeaderMenuBtn = $('#document-info-bar-item-menu');
const documentHeaderMenuShareLinkCreate = $("#document-header-menu-item-share-link-create")
const documentHeaderMenuShareLinkTextArea = $("#document-header-menu-item-share-link")
const documentHeaderMenuShareLinkCopy = $("#document-header-menu-item-share-link-copy")
const documentHeaderMenuDeleteDocument = $("#document-header-menu-item-delete-document")
const titleElement = $("#document-info-bar-item-title");
let menuOpen = false;
var userMap = new Map();
var readOnlyMessage ="（只读）"

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
            documentHeaderMenuShareLinkTextArea.val(window.location.protocol + "//" + window.location.host+"/share/" + link);
        });
    }
}
const deleteDocument = () => {
    const confirmDelete = confirm("确认删除此文档?")
    if (confirmDelete){
        $.ajax({
            url: `/api/user/${pageData.userId}/document/${pageData.documentId}/delete`,
            type: 'POST',
            timeout: 3e4,
        }).done((data) => {
          if (data.status === 1) {
              alert("文档删除成功");
              location.href = "/";
          } else{
              alert("文档删除失败");
          }
        })
    }
}
const copyDocumentShareLink = () => {
    documentHeaderMenuShareLinkTextArea.select();
    document.execCommand("copy");

}
const getAllDocumentUsers = () => {
    $.ajax({})
}
const changeTitle = ()=>{
    const textArea = $("<textarea class='title-textarea' id='title-textarea'></textarea>");
    titleElement.text("");
    titleElement.append(textArea);
    textArea.val(pageData.title);

    textArea.on("keydown", (event)=>{
        if (event.keyCode === 13) {
            titleElement.text(textArea.val());
        }
    })
}
const loadDocumentUsers = (data) => {
    const users = data.documentUsers;
    let usersString = ""
    //let userInfo = []
    for (let i = 0; i < users.length; i++) {
        const user = users[i];
        const userId = user.user_id;
        usersString += userId + ";";
    }
    usersString = usersString.substring(0, usersString.length - 1)
    $.ajax({
        url: `/api/user/${pageData.userId}/info/${usersString}`,
        type: 'GET',
        timeout: 3e4,

    }).done((data) => {
        const users = data;
        for (let i = 0; i < users.length; i++) {
            const user = users[i];
            const userId = user.user_id;
            userMap.set(userId, user)
        }
    }).done(() => {
        for (let i = 0; i < users.length; i++) {
            const user = users[i];
            const userId = user.user_id;
            const userName = userMap.get(userId).user_name;
            const permissionType = user.permission_type;
            let permissionItem = $(`<ul>${userName}:<input type="radio" class="permission-radio permission-${userId}" id="editable_${userId}" name="permission_${userId}" value="true">可编辑  <input type="radio" class="permission-radio permission-${userId}" id="readonly_${userId}" name="permission_${userId}" value="false">只读</ul>`);
            $("#document-header-menu-item-share-link-users").append(permissionItem)
            if (permissionType) {
                $(`#editable_${userId}`).prop("checked",true)
            } else {
                $(`#readonly_${userId}`).prop("checked",true)
            }
        }
        //$(`#editable-${userMap.get(pageData.userId)}`).prop("disabled",true);
        if(pageData.authorId !== pageData.userId){
            $(`.permission-radio`).prop("disabled",true);
        } else{
            $(`.permission-radio`).on("change",changeUserPermission);
        }
    })
}
const changeUserPermission = (event) => {
    const targetUserId = event.target.id.split("_")[1]
    $.ajax({
        url:`/api/user/${pageData.userId}/document/${pageData.documentId}/permission/${targetUserId}/${$(event.target).val()}`,
        type:'POST',
        timeout:3e4,
    }).done((data)=>{
        if(data.status !== 1){
            alert("权限修改失败");
        }
    });
}

documentHeaderMenuShareLinkCreate.on('click', getDocumentShareLink)
documentHeaderMenuBtn.on('click', () => {
    changeMenuState();
});
documentHeaderMenuShareLinkCopy.on('click', copyDocumentShareLink)
documentHeaderMenuDeleteDocument.on('click', deleteDocument)
if (pageData.permissionType){
    titleElement.on("dblclick",changeTitle)
}
history.pushState({}, '', `/document/${pageData.documentId}`)