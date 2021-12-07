$(document).ready(function () {
    sortTable($(".customTable"), 0)
});

function sortTable(table, columnIndexSort) {
    console.log(table.find("tbody tr"))
    swap(table, 0, 3)
}

function quickSortTableRec(table, iStart, iEnd) {

}

function swap(table, i1, i2) {
    tbody = table.find("tbody")
    elt1 = tbody.children().get(i1)
    afterElt1 = elt1.next
    elt2 = tbody.children().get(i2)
    tbody.children().insertBefore(elt1, elt2)
}