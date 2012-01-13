// copied from
// http://www.fletchzone.com/post/jQuery-Unobtrusively-Animated-Add-and-Remove-Table-Rows.aspx
(function ($) {
    var defaults = {
        rowSpeed: 300,
        newRow: null,
        addTop: true,
        removeTop: true
    };
    var newClasses = "newRow"
    var options = $.extend(defaults, options);
    $.fn.addRow = function (options) {
        opts = $.extend(defaults, options);
        var $table = $(this);
        var $tableBody = $("tbody", $table);
        var t = $(opts.newRow).find("td").wrapInner("<div style='display:none;'/>").parent()
        if (opts.addTop) t.appendTo($tableBody);
        else t.prependTo($tableBody);
        t.attr("class", newClasses).removeAttr("id").show().find("td div").slideDown(options.rowSpeed, function () {
            $(this).each(function () {
                var $set = jQuery(this);
                $set.replaceWith($set.contents());
            }).end()
        })
        return false;
    };
    $.fn.removeRow = function (options) {
        opts = $.extend(defaults, options);
        var $table = $(this);
        var t
        if (opts.removeTop) t = $table.find('tbody tr:last')
        else t = $table.find('tbody tr:first');
        t.find("td")
        .wrapInner("<div  style='DISPLAY: block'/>")
        .parent().find("td div")
        .slideUp(opts.rowSpeed, function () {
            $(this).parent().parent().remove();
        });
        return false;
    };
    return this;
})(jQuery);
