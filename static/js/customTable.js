$(document).ready(function () {
    // Make initial sorting
    sortTable($(".customTable.sorted"))

    // Add buttons on the table header to sort it as wanted
    $(".customTable thead tr th").click(function (e) {
        e.preventDefault();
        newSortColumnIndex = $(this).index()
        table = $(this).parents(".customTable")[0]
        if (parseInt(table.style.getPropertyValue('--sortColumnIndex')) == newSortColumnIndex) { // This part needs to be optimised (don't sort just reverse elements order)
            if (table.style.getPropertyValue('--sortOrder') == 0) {
                table.style.setProperty('--sortOrder', 1)
            } else {
                table.style.setProperty('--sortOrder', 0)
            }
        } else {
            table.style.setProperty('--sortColumnIndex', newSortColumnIndex)
            table.style.setProperty('--sortOrder', 0)
        }
        sortTable($(this).parents(".customTable"))
    });
});

function sortTable(table) {
    tbody = table.find("tbody")[0]
    if (tbody != undefined) {
        trList = tbody.getElementsByTagName("tr")
        
        sortColumnIndex = parseInt(table[0].style.getPropertyValue('--sortColumnIndex'))
        quickSortTrList(
            trList, 
            0, 
            tbody.children.length - 1, 
            sortColumnIndex, 
            parseInt(table[0].style.getPropertyValue('--sortOrder')),
            parseInt(table.parent().find("thead")[0].children[0].children[sortColumnIndex].style.getPropertyValue('--dataType'))
            )
    }
}

function quickSortTrList(trList, iStart, iEnd, columnIndexSort, sortOrder, columnDataType) {
    if (iStart < iEnd) {
        pi = quickSortTrListPartition(trList, iStart, iEnd, columnIndexSort, sortOrder, columnDataType)
        quickSortTrList(trList, iStart, pi - 1, columnIndexSort, sortOrder, columnDataType)
        quickSortTrList(trList, pi + 1, iEnd, columnIndexSort, sortOrder, columnDataType)
    }
}

function quickSortTrListPartition(trList, iStart, iEnd, columnIndexSort, sortOrder, columnDataType) {
    pivot = parseIfNeeded(trList[iEnd].children[columnIndexSort].innerHTML, columnDataType)

    i = iStart - 1

    for (j = iStart; j <= iEnd - 1; j++) {
        value = parseIfNeeded(trList[j].children[columnIndexSort].innerHTML, columnDataType)
        if (quickSortVerifyOrder(value, pivot, sortOrder)) {
            i++
            swap(trList, i, j)
        }
    }
    swap(trList, i + 1, iEnd)
    return i + 1
}

function parseIfNeeded(value, dataType) {
    if (dataType == 1) { // If the value is a number we remove spaces, balises and parse it
        value = value.replace(/<span>/, '')
        value = value.replace(/<\/span>/, '')
        return parseInt(value.replace(/ /g, ''))
    }
    return value
}

// order == 1 means <
// else means >
function quickSortVerifyOrder(value, pivot, order) {
    if (order == 1) {
        return value < pivot
    } else {
        return value > pivot
    }
}

function swap(trList, i1, i2) {
    if (i1 != i2) {
        element1 = trList[i1]
        element2 = trList[i2]

        var clonedElement1 = element1.cloneNode(true);
        var clonedElement2 = element2.cloneNode(true);

        element2.parentNode.replaceChild(clonedElement1, element2);
        element1.parentNode.replaceChild(clonedElement2, element1);
    }
}