const markdownText = $("#document-content");
const markdownViewer = $("#document-viewer-container");
const favouriteBtn = $("#document-info-bar-item-favourite");
const saveBtn = $("#document-info-bar-item-save");
const refreshBtn = $("#document-info-bar-item-refresh");
const currentTitle = $("#document-info-bar-item-title");
const downloadBtn=$("#document-header-menu-item-download");
const autoSaveCheckbox = $("#document-header-menu-item-auto-save")


const viewportHeight = window.innerHeight;
const viewportWidth = window.innerWidth;
const initialRowHeight = viewportHeight / 30;
const initialColumnWidth = 5 * initialRowHeight;
const excelColumnHeaders = $("#excel-column-headers")
const excelTable = $("#excel-table")


let chosenCellList = []
let startCell = null;
let endCell = null;
let startRow = null;
let endRow = null;
let startColumn = null;
let endColumn = null;
let initialSuccessful = false;
let excelData = [[]];
let isMouseDown = false;
let lastRefreshTime = new Date().getTime();
let refreshTime = 3000;
let refreshCount = 0;

const initializeExcelTable = () => {
    const initialTableRows = 30;
    const initialTableColumns = viewportWidth / initialColumnWidth;
    for (let i = 0; i < initialTableColumns; i++) {
        const columnHeaderTitle = columnNumberToLetter(i);
        excelColumnHeaders.append(`<th class="excel-column-header excel-column-${i}"><div class="excel-column-header-cell excel-row-0 excel-column-${i}">${columnHeaderTitle}</div></th>`);
    }
    for (let j = 1; j < initialTableRows; j++) {
        excelTable.append(`<tr class="excel-data excel-row-${j}"></tr>`);
        const excelRow = $(`.excel-row-${j}`);
        excelRow.append(`<td class="excel-row-header excel-column-0"><div class="excel-row-header-cell excel-row-${j} excel-column-0">${j}</div></td>`);

        for (let i = 1; i < initialTableColumns; i++) {

            excelRow.append(`<td class="excel-data-cell excel-column-${i}"><div class="excel-data-cell excel-row-${j} excel-column-${i}"></div></td>`);
        }

    }
    const allCells = $(".excel-data-cell, .excel-column-header-cell, .excel-row-header-cell");
    allCells.css("height", initialRowHeight);
    allCells.css("width", initialColumnWidth);
}

const getExcelData = () => {
    $.ajax({
        url: `/api/user/${pageData.userId}/document/${pageData.documentId}/get`,
        type: "GET",
        timeout: 3e4,
        dataType: "json",
    }).done(function (response) {
        loadDocumentUsers(response);
        const content = response.content;
        const excelCells = JSON.parse(content).excelCells;
        for (const row in excelCells) {
            for (const column in excelCells[row]) {
                const cell = $(`.excel-row-${row} .excel-column-${column} .excel-data-cell`);
                cell.text(excelCells[row][column].content);
                cell.css(excelCells[row][column].style)
            }
        }
        initialSuccessful = true;
        autoRefreshInterval = setInterval(refreshExcelData, refreshTime);
    });
}
const saveExcelDocument = () => {
    let postData;
    if (currentTitle.text() === pageData.title) {
        postData = {}
    } else {
        postData = {
            "title": currentTitle.text(),

        }
    }
    $.ajax({
        type: "POST",
        url: `/api/user/${pageData.user_id}/document/${pageData.document_id}/save`,
        dataType: "application/json",
        data: JSON.stringify(postData),
    })
}
const removeUnableCells = () => {
    const currentTime = Math.floor(Date.now() / 1000)
    const unableCells = $(".unable-edit-cell");
    for (let i = 0; i < unableCells.length; i++) {
        const unableCell = $(unableCells[i]);
        const addTime = unableCell.data("add-time");
        if (currentTime - addTime > 20) {
            unableCell.find(".other-user-updated-cell").remove();
            unableCell.removeClass("unable-edit-cell");
        }

    }
}
const fetchUserData= (userId,usersString)=>{
    return $.ajax({
        url: `/api/user/${userId}/info/${usersString}`,
        type: 'GET',
        timeout: 30000
    });
}
const refreshExcelData = () => {

    removeUnableCells();
    let postData = {
        "timestamp": lastRefreshTime
    }
    postData = JSON.stringify(postData)
    $.ajax({
        url: `/api/user/${pageData.userId}/document/${pageData.documentId}/excel/refresh`,
        type: "POST",
        timeout: 3e4,
        dataType: "json",
        data: postData,
    }).done(async (data) => {
        lastRefreshTime = new Date().getTime();
        const updateCells = data.cells;

        let UpdatedFrom = new Map();
        for (let i = 0; i < updateCells.length; i++) {
            const userId = updateCells[i].userId;
            if (userId === pageData.userId) continue;
            const row = updateCells[i].ReceivedExcelCell.row;
            const column = updateCells[i].ReceivedExcelCell.column;
            const content = updateCells[i].ReceivedExcelCell.content;
            const style = updateCells[i].ReceivedExcelCell.style;
            const cell = $(`.excel-row-${row} .excel-column-${column} .excel-data-cell`);

            UpdatedFrom.set(userId, cell)

            cell.text(content);
            cell.css(style);
        }
        if (UpdatedFrom.size === 0) {
            refreshCount++;
            if (refreshCount > 6) {
                refreshTime = 10000;
                clearInterval(autoRefreshInterval)
                autoRefreshInterval = setInterval(refreshExcelData, refreshTime)
            }
            return
        } else {
            refreshCount = 0;
            refreshTime = 3000;
            clearInterval(autoRefreshInterval)

        }
        for (let [userId, cell] of UpdatedFrom) {
            let userName = "";
            let user = userMap.get(userId)
            if (user === undefined) {
                const users = await fetchUserData(pageData.userId, userId);
                const user = users[0];
                userName = user.user_name;
                userMap.set(userId, user)
            }
            userName = user.user_name;
            const cellParent = cell.parent();
            cellParent.addClass("unable-edit-cell");
            cellParent.data("add-time", Math.floor(Date.now() / 1000));
            $(`.user-${userId}-cell`).remove();
            const OtherUserCell = $(`<div class="other-user-updated-cell user-${userId}-cell"><label class="other-user-updated-cell-label">${userName}</label><div class="other-user-updated-cell-border"></div></div>`)
            OtherUserCell.css({
                height: (cell.height() + 2) + "px",
                width: (cell.width() + 2) + "px",
                left: "-1px",
                top: "-1px"
            });
            cellParent.append(OtherUserCell);
        }
        autoRefreshInterval = setInterval(refreshExcelData, refreshTime)

    })
}
const columnNumberToLetter = (column) => {
    let columnLetter = '';
    while (column > 0) {

        let remainder = (column - 1) % 26;
        columnLetter = String.fromCharCode(65 + remainder) + columnLetter;
        column = Math.floor((column - remainder) / 26);
    }
    return columnLetter;
}


const GetCellRowAndColumn = ($cell) => {
    const classList = $cell.attr("class");
    const rowRegex = /excel-row-(\d+)/;
    const colRegex = /excel-column-(\d+)/;
    const rowMatch = classList.match(rowRegex);
    const rowNumber = rowMatch ? parseInt(rowMatch[1]) : null;
    const colMatch = classList.match(colRegex);
    const colNumber = colMatch ? parseInt(colMatch[1]) : null;
    return {row: rowNumber, column: colNumber};
}
const changeSelectedCellsStyle = (startRow, startColumn, endRow, endColumn) => {
    $(".chosen-cell").removeClass("chosen-cell")
    $(".chosen-cell-top").removeClass("chosen-cell-top")
    $(".chosen-cell-bottom").removeClass("chosen-cell-bottom")
    $(".chosen-cell-left").removeClass("chosen-cell-left")
    $(".chosen-cell-right").removeClass("chosen-cell-right")

    let minRow = Math.min(startRow, endRow);
    let maxRow = Math.max(startRow, endRow);
    let minColumn = Math.min(startColumn, endColumn);
    let maxColumn = Math.max(startColumn, endColumn);
    for (let i = minRow; i <= maxRow; i++) {
        for (let j = minColumn; j <= maxColumn; j++) {
            $(`.excel-row-${i} .excel-column-${j}`).addClass("chosen-cell")
            $(`.excel-row-${minRow} .excel-column-${j} .chosen-cell`).addClass("chosen-cell-top")
            $(`.excel-row-${maxRow} .excel-column-${j} .chosen-cell`).addClass("chosen-cell-bottom")
            $(`.excel-row-${i} .excel-column-${minColumn} .chosen-cell`).addClass("chosen-cell-left")
            $(`.excel-row-${i} .excel-column-${maxColumn} .chosen-cell`).addClass("chosen-cell-right")

        }
    }
}
const detectMouseDown = (event) => {
    const unable = event.target.parentNode.classList.contains("unable-edit-cell");
    if (unable) return;
    removeTextareaAndSaveValue();
    if (!isMouseDown) {
        $(".chosen-cell").removeClass("chosen-cell")
        $(".chosen-cell-top").removeClass("chosen-cell-top")
        $(".chosen-cell-bottom").removeClass("chosen-cell-bottom")
        $(".chosen-cell-left").removeClass("chosen-cell-left")
        $(".chosen-cell-right").removeClass("chosen-cell-right")
        $(".chosen-cell-selected").removeClass("chosen-cell-selected")
        startCell = event.target;
        isMouseDown = true;
        const startCellRowAndColumn = GetCellRowAndColumn($(startCell));
        startRow = startCellRowAndColumn.row;
        startColumn = startCellRowAndColumn.column;
    }
}
const detectMouseOver = (event) => {
    if (isMouseDown) {

        $(".chosen-cell-selected").removeClass("chosen-cell-selected")
        const selectedCell = $(event.target);
        selectedCell.addClass("chosen-cell-selected")
        const selectedCellRowAndColumn = GetCellRowAndColumn(selectedCell);
        endRow = selectedCellRowAndColumn.row;
        endColumn = selectedCellRowAndColumn.column;
        changeSelectedCellsStyle(startRow, startColumn, endRow, endColumn)
    }
}
const detectMouseUp = (event) => {

    endCell = event.target;
    $(endCell).addClass("chosen-cell-selected")
    isMouseDown = false;
}
const detectDoubleClick = (event) => {
    const unable = event.target.parentNode.classList.contains("unable-edit-cell");
    if (unable) return;
    removeTextareaAndSaveValue();
    const doubleClickedCell = $(event.target);
    const textareaElement = $("<textarea class='cell-textarea'></textarea>");
    // const doubleClickedCellOffset = doubleClickedCell.offset();
    // const doubleClickedCellHeight = doubleClickedCell.height();
    // const doubleClickedCellWidth = doubleClickedCell.width();

    doubleClickedCell.append(textareaElement);
    textareaElement.focus();
}
const removeTextareaAndSaveValue = () => {
    const textareaElement = $(".cell-textarea");
    const cellValue = textareaElement.val();
    const cellRowAndColumn = GetCellRowAndColumn($(endCell));
    const cellRow = cellRowAndColumn.row;
    const cellColumn = cellRowAndColumn.column;
    textareaElement.remove();
    if (cellValue !== undefined) {
        const postData = {
            "row": cellRow,
            "column": cellColumn,
            "content": cellValue,
            "style": {}
        }
        $.ajax({
            url: `/api/user/${pageData.userId}/document/${pageData.documentId}/excel/update`,
            type: "POST",
            timeout: 3e4,
            dataType: "json",
            data: JSON.stringify(postData),
        }).done(function (response) {
            if (response.status === 1) {
                const cellTobeUpdated = $(`.excel-row-${cellRow} .excel-column-${cellColumn} .excel-data-cell`);
                cellTobeUpdated.text(cellValue);
            }

        })
    }
}


initializeExcelTable();
getExcelData();


excelTable.on("mousedown",  (event)=> {
    detectMouseDown(event);
});

excelTable.on("mouseover", (event) =>{
    detectMouseOver(event);
});
excelTable.on("mouseup",(event)=> {
    detectMouseUp(event);
});
if (pageData.permissionType) {
    excelTable.on("dblclick",  (event)=> {
        detectDoubleClick(event);
    });
    saveBtn.on("click", saveExcelDocument);

} else{
    currentTitle.text(currentTitle.text()+readOnlyMessage);
}

excelTable.on("keydown",  (event)=> {
    if (event.keyCode === 13) {
        removeTextareaAndSaveValue();
    }
});

refreshBtn.on("click", refreshExcelData);
downloadBtn.on("click", ()=> {
    $.ajax({
        url: `/api/user/${pageData.userId}/document/${pageData.documentId}/excel/download`,
        type:`GET`,
        xhrFields: {
            responseType: 'blob'
        },
    }).done((response, status, xhr)=>{
       const blob = new Blob([response], { type: xhr.getResponseHeader('Content-Type') });
        const url = window.URL.createObjectURL(blob);


        const a = document.createElement('a');
        a.href = url;
        a.download = xhr.getResponseHeader('Content-Disposition').split('filename=')[1];
        document.body.appendChild(a);
        a.click();

        window.URL.revokeObjectURL(url);
        document.body.removeChild(a);
    })
})
let autoRefreshInterval = 0;
