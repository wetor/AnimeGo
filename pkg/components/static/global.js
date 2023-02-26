$(function () {
    $('[data-toggle="popover"]').popover()
})

// 初始化列表
$(".list-group").each(function() {
    reorderItems($(this));
    refreshButtonStatus($(this));
});

// 添加新项
$(".a-array").on("click", ".a-list-adder", function() {
    var $listGroup = $(this).closest(".a-array").find(".list-group");
    var $newItem = createNewItem($listGroup.find(".a-template"));
    $listGroup.find('li:last').before($newItem);
    $newItem.fadeOut(150).fadeIn(300);
    reorderItems($listGroup);
    refreshButtonStatus($listGroup);
    $('[data-toggle="popover"]').popover();
});

// 删除当前项
$(".a-array").on("click", ".a-list-deleter", function() {
    var $listItem = $(this).closest("li");
    $listItem.remove();
    reorderItems($listItem.closest(".list-group"));
    refreshButtonStatus($listItem.closest(".list-group"));
});

// 上移当前项
$(".a-array").on("click", ".a-list-up", function() {
    var $listItem = $(this).closest("li");
    var $prevItem = $listItem.prev();
    if ($prevItem.length && !$prevItem.hasClass("a-template")) {
        $listItem.insertBefore($prevItem);
        $listItem.fadeOut(150).fadeIn(300);
        $listItem.find('.a-item-body').collapse('show');
        $prevItem.find('.a-item-body').collapse('hide');
        reorderItems($listItem.closest(".list-group"));
        refreshButtonStatus($listItem.closest(".list-group"));
    }
});

// 下移当前项
$(".a-array").on("click", ".a-list-down", function() {
    var $listItem = $(this).closest("li");
    var $nextItem = $listItem.next();
    if ($nextItem.length && !$nextItem.hasClass("a-template")) {
        $listItem.insertAfter($nextItem);
        $listItem.fadeOut(150).fadeIn(300);
        $listItem.find('.a-item-body').collapse('show');
        $nextItem.find('.a-item-body').collapse('hide');
        reorderItems($listItem.closest(".list-group"));
        refreshButtonStatus($listItem.closest(".list-group"));
    }
});

function refreshButtonStatus($listGroup) {
    $listGroup.find("li:not(.a-template)").each(function(index) {
        var $listItem = $(this);
        var $prevItem = $listItem.prev();
        var $nextItem = $listItem.next();

        $listItem.find(".a-list-up").prop("disabled", !$prevItem.length || $prevItem.hasClass("a-template"));
        $listItem.find(".a-list-down").prop("disabled", !$nextItem.length || $nextItem.hasClass("a-template"));
    });
}

// 重新排列项的序号
function reorderItems($listGroup) {
    var itemNumber = 0;
    var name = $listGroup.closest(".a-array").attr('data-name') + '__list__'
    $listGroup.find(".list-group-item:not(.a-template)").each(function() {
        var $header = $(this).find('.a-item-header div')
        var $collapse = $(this).find('.a-item-body')
        $header.find('span').text("[" + itemNumber + "]");
        var name_index = name + '-' + itemNumber
        $header.attr('data-target', '#item-' + name_index)
        $collapse.attr('id', 'item-' + name_index)
        $(this).find('.a-input').each(function(i, e){
            last = $(e).closest('.form-group').attr('data-name')
            $(e).attr('id', name_index + '-' + last)
            $(e).attr('name', name_index + '-' + last)
        })

        itemNumber++;
    });
}

// 克隆一个新的list-group-item，并修改其中的文本
function createNewItem($template) {
    var $newItem = $template.clone();
    $newItem.removeClass("a-template");
    $newItem.removeAttr('hidden');
    return $newItem;
}

function traverse(obj) {
    for (let key in obj) {
        if (typeof obj[key] === 'object') {
            if(key.endsWith('__list__')){
                var new_key = key.slice(0, -'__list__'.length)
                obj[new_key] = Object.values(obj[key]);
                delete obj[key]
                traverse(obj[new_key]);
            }else{
                traverse(obj[key]);
            }
        }
    }
}

$(document).ready(function() {
    $('#myForm').submit(function(event) {
        event.preventDefault();
        var formData = {};

        // Loop through all form elements
        $(this).find(':input:not(.a-template *,:submit,:button)').each(function() {
            var val = $(this).val();
            var name = $(this).attr('name');
            var keys = name.split('-');
            var target = formData;
            for (var i = 0; i < keys.length - 1; i++) {
                var key = keys[i];
                if (!target[key]) {
                    target[key] = {};
                }
                target = target[key];
            }
            var finalKey = keys[keys.length - 1];
            if ($.isNumeric(val)) {
                val = parseFloat(val);
            } else if (val === 'true' || val === 'false') {
                val = (val === 'true');
            }
            if (target[finalKey]) {
                if ($.isArray(target[finalKey])) {
                    target[finalKey].push(val);
                } else {
                    target[finalKey] = [target[finalKey], val];
                }
            } else {
                target[finalKey] = val;
            }
        });

        traverse(formData);

        // Convert object to JSON and submit the form
        var jsonData = JSON.stringify(formData.Config);
        console.log(jsonData)
        $.ajax({
            type: 'POST',
            url: 'test',
            data: { data: jsonData },
            success: function() {
                console.log('Form submitted successfully!');
            },
            error: function() {
                console.log('Error submitting form.');
            }
        });
    });
});
