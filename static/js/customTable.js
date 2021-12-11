$(document).ready(function () {
    tableToSort = $(".customTable.sorted")
    sortTable(tableToSort, parseInt(tableToSort[0].style.getPropertyValue('--sortColumnIndex')), parseInt(tableToSort[0].style.getPropertyValue('--sortOrder'))) // Uncomment this last parameter to sort in ascending order
});

function sortTable(table, columnIndexSort, sortOrder) {
    tbody = table.find("tbody")[0]
    if (tbody != undefined) {
        trList = tbody.getElementsByTagName("tr")
        /*for (i = 0; i < trList.length; i++) {
            console.log(trList[i].children[columnIndexSort].innerHTML + " - " + trList[i].children[columnIndexSort + 1].innerHTML)
        }*/
        quickSortTrList(trList, 0, tbody.children.length-1, columnIndexSort, sortOrder)
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
    pivot = parseInt(trList[iEnd].children[columnIndexSort].innerHTML)
    //console.log(pivot + " - " + trList[iEnd].children[1].innerHTML)
    i = iStart - 1

    for (j = iStart; j <= iEnd - 1; j++) {
        if (quickSortVerifyOrder(parseInt(trList[j].children[columnIndexSort].innerHTML), pivot, sortOrder)) {
            i++
            swap(trList, i, j)
        }
    }
    swap(trList, i + 1, iEnd)
    return i + 1
}

// order == 1 means <
// else means >
function quickSortVerifyOrder(value1, value2, order) {
    if (order == 1) {
        return value1 < value2
    } else {
        return value1 > value2
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