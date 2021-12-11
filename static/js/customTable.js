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
        /*for (i = 0; i < trList.length; i++) {
            console.log(trList[i].children[columnIndexSort].innerHTML + " - " + trList[i].children[columnIndexSort + 1].innerHTML)
        }*/
        quickSortTrList(trList, 0, tbody.children.length - 1, parseInt(table[0].style.getPropertyValue('--sortColumnIndex')), parseInt(table[0].style.getPropertyValue('--sortOrder')))
    }
}

function quickSortTrList(trList, iStart, iEnd, columnIndexSort, sortOrder) {
    if (iStart < iEnd) {
        pi = quickSortTrListPartition(trList, iStart, iEnd, columnIndexSort, sortOrder)
        quickSortTrList(trList, iStart, pi - 1, columnIndexSort, sortOrder)
        quickSortTrList(trList, pi + 1, iEnd, columnIndexSort, sortOrder)
    }
}

function quickSortTrListPartition(trList, iStart, iEnd, columnIndexSort, sortOrder) {
    pivot = parseIfNeeded(trList[iEnd].children[columnIndexSort].innerHTML)

    //console.log(pivot + " - " + trList[iEnd].children[1].innerHTML)
    i = iStart - 1

    for (j = iStart; j <= iEnd - 1; j++) {
        value = parseIfNeeded(trList[j].children[columnIndexSort].innerHTML)
        if (quickSortVerifyOrder(value, pivot, sortOrder)) {
            i++
            swap(trList, i, j)
        }
    }
    swap(trList, i + 1, iEnd)
    return i + 1
}

function parseIfNeeded(value) {
    if (isNaN(value) || value > 2000) { // Ugly fix but xd
        return value.replace(/ /g, '')
    }
    return parseInt(value)
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