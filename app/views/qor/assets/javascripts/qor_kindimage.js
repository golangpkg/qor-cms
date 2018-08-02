(function (factory) {
    if (typeof define === 'function' && define.amd) {
        define(['jquery'], factory);
    } else if (typeof exports === 'object') {
        factory(require('jquery'));
    } else {
        factory(jQuery);
    }
})(function ($) {
    'use strict';
    let componentHandler = window.componentHandler,
        NAMESPACE = 'qor.kindimage',
        EVENT_ENABLE = 'enable.' + NAMESPACE,
        EVENT_DISABLE = 'disable.' + NAMESPACE,
        EVENT_UPDATE = 'update.' + NAMESPACE,
        SELECTOR_COMPONENT = '[class*="mdl-js"],[class*="mdl-tooltip"]';

    function enable(target) {
        /*jshint undef:false */
        var uploadEditor = KindEditor.editor({
            uploadJson: '/admin/common/kindeditor/upload?dir=image',
            allowFileManager: false
        });
        $('#imgUrlButton').click(function () {
            uploadEditor.loadPlugin('uploadImage', function () {
                uploadEditor.plugin.fileDialog({
                    clickFn: function (data) {
                        $("#imgUrlPreview").attr("src", data.url);
                        $("#imgUrlHidden").val(data.url);
                        uploadEditor.hideDialog();
                    }
                });
            });
        });
    }

    function disable(target) {
    }
    $(function () {
        $(document)
            .on(EVENT_ENABLE, function (e) {
                enable(e.target);
            })
            .on(EVENT_DISABLE, function (e) {
                disable(e.target);
            })
            .on(EVENT_UPDATE, function (e) {
                disable(e.target);
                enable(e.target);
            });
    });
});