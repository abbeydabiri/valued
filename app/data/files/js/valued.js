// All JS Triggers


 // function to possibly override keypress
trapBackspace = function(event) {
	var keynum;
	if (window.event) {
		keynum = window.event.keyCode;
	} else if (event.which) {
		keynum = event.which;
	}

	if (keynum == 8) { 
		return false;
	}
	return true;
}



function appSubscribeOption(idTag) {
	$("input[type=radio]").each(function(index) { 
		$(this).prop("checked", false);
		if ($(this).prop("value") == idTag) {
			$(this).prop("checked", true);  
		}
	});
}


function slideBarClicked(idTag) {
	$(".slidebarToggle")
		.each(function(index) { $(this).removeClass('active'); });
	$(idTag).addClass('active');


	var filter = $(idTag).attr('filter');
	var page = $('#pageId').val();
}

function slideBarLeft() {
	return;
	$('#slidebarToggle1').css('display', 'inline-block');
	// $('#slidebarToggle2').css('display','inline-block');
	$('#slidebarToggle3').css('display', 'inline-block');
	$('#slidebarToggle4').css('display', 'inline-block');
	$('#slidebarToggle5').css('display', 'none');
}

function slideBarRight() {
	return;
	$('#slidebarToggle1').css('display', 'none');
	// $('#slidebarToggle2').css('display','none');
	$('#slidebarToggle3').css('display', 'inline-block');
	$('#slidebarToggle4').css('display', 'inline-block');
	$('#slidebarToggle5').css('display', 'inline-block');
}

function addToFavorites(idTag, page, signedIn) {
	if (signedIn !== 'yes') {

		pageCamelcase = "";
		switch(page) {
			case "reward":
			pageCamelcase = "Rewards"; 
			break;
			case "merchant":
			pageCamelcase = "Merchants"; 
			break;
		}
		warningMessage('Please LOGIN to add '+pageCamelcase+' to your favourites');
		return
	}

	var src = $(idTag).attr("src");
	var title = encodeURIComponent($(idTag).attr("title"));
	var control = $(idTag).attr("control");
	urlPost = "/app-favorite?" + page + "control=" + control + "&title=" +
			  title + "&action=";


	if (src.match("heart_filled.png$")) {
		$(idTag).attr("src", "files/img/heart.png");
		quickForm(urlPost + "remove");
	} else {
		$(idTag).attr("src", "files/img/heart_filled.png");
		quickForm(urlPost + "add");
	}
}

function toggleAppSidebar(idTag) {
	$('#' + idTag).toggle('slide', {direction: "left"}, 0);
}


function defaultImage(className) {
	$(className).each(function(pos, imageBox) {
		var imageboxIMG = imageBox.getElementsByTagName("img")[0];
		if ($(imageboxIMG).attr('src') == "") {
			$(imageboxIMG).attr('src', 'files/img/default.jpg');
		}
	});
}


function editFormCompressImage(imageID) {
	var imageInputFile = $('#' + imageID + 'File').get(0).files;

	if (imageInputFile.length == 0) {
		$.notify({message: "<center>Please Select File(s) to Upload</center>"},
				 {type: 'danger', timer: 5000});
		return;
	}

	if (imageInputFile.length > 1) {
		$.notify({message: "<center>Exceeded Maximum Allowed Files</center>"},
				 {type: 'danger', timer: 5000});
		return;
	}
	

	for (var counter = 0; counter < imageInputFile.length; counter++) {
		var reader = new FileReader();
		reader.onload = function(file_name) {
			return function(event) {
				var the_url = event.target.result;

				if (the_url.indexOf("base64") == 5) {
					mime_type = "image/jpeg";
					the_url = the_url.replace("data:base64","data:image/jpeg;base64");
				}
				
				var mime_type = compressImage(file_name, the_url)
				if (mime_type != "") {
					$('#' + imageID + 'Src').css('height',"30px !important");
					$('#' + imageID + 'Src').attr('src', '../files/img/siteloader.gif');
				}


				var source_img_obj = document.createElement('img');
				source_img_obj.src = the_url;

				var callback = function(source) {
					if(!source) source = this;

					var quality = 80;
					var maxWidth = 305;
					maxWidth = maxWidth || 1000;
					var natW = source_img_obj.naturalWidth;
					var natH = source_img_obj.naturalHeight;
					var ratio = natH / natW;
					
					natW = maxWidth;
					natH = ratio * maxWidth;


					cvs = document.createElement('canvas');
					cvs.width = natW;
					cvs.height = natH;


					var ctx = cvs.getContext("2d").drawImage(source_img_obj, 0, 0, natW, natH);
					var compressed_image_url = cvs.toDataURL(mime_type, quality / 100);
					
					$('#' + imageID).val(compressed_image_url);
					$('#' + imageID + 'Src').attr('src', compressed_image_url);
					$('#' + imageID + 'Src').css('height',"100% !important");
					$('#' + imageID + 'Name').val(file_name);


				}

				if (mime_type != "") {
					source_img_obj.onload = callback;
				}

			};
		}(imageInputFile[counter].name);
		reader.readAsDataURL(imageInputFile[counter]);
	}
}

function tinyMceAddImages() {
	var imageBoxHtml = "";
	$('.imagebox')
		.each(function(pos, imageBox) {

			var checkBox = imageBox.getElementsByTagName("input")[0];
			if (checkBox.checked) {
				imageBox.removeChild(checkBox);

				var anchorTag = imageBox.getElementsByTagName("a")[0];
				var imageCaption =
					anchorTag.getElementsByClassName("imagecaption")[0];
				anchorTag.removeChild(imageCaption);

				$(imageBox).attr(
					"style",
					"display: inline-table !important; position: relative; border: 1px solid #ddd; border-radius: 4px; -webkit-transition: all .2s ease-in-out; margin: 3px;")

					imageBoxHtml += imageBox.outerHTML;
			}
		});

	tinyMCE.activeEditor.execCommand('mceInsertContent', false, imageBoxHtml);
}

$('html').on("submit", "#searchMediaLibrary", function(event) {
	event.preventDefault();
	var formData = new FormData($(this)[0]);
	var formUrl = $('#searchMediaLibrary').attr('action') + "?";
	submitForm(formUrl, formData);
	return false;
});


function compressImage(file_name, the_url) {
	var mime_type = "";

	if (the_url.indexOf("image/jpeg;base64") == 5) {
		mime_type = "image/jpeg";
	}

	if (the_url.indexOf("image/png;base64") == 5) {
		mime_type = "image/png";
	}

	if (the_url.indexOf("image/gif;base64") == 5) {
		mime_type = "image/gif";
	}

	if (mime_type.length == 0) {
		$.notify(
			{
			  message:
				  "<center>File (<b>" + file_name + "</b>) Not Allowed</center>"
			},
			{type: 'danger', timer: 5000});
	}

	return mime_type;
}


var delay = (function() {
	var timer = 0;
	return function(callback, ms) {
		clearTimeout(timer);
		timer = setTimeout(callback, ms);
	};
})();

function adjustFooter(){
	delay(function() {
		appHeader = $('.appHeader').outerHeight();
		appFooter = $('.appFooter').outerHeight();
		$('.appContent').css('height','calc(100% - '+(appFooter+appHeader)+'px)');
		$('#searchFilterDrop').css('max-height',$('.appContent').css('height'));
	}, 500);
}

function setupDropdownSearch(url, idTag) {
	if (url == '') {
		url = idTag;
	}


	$('#' + idTag + 'title')
		.keyup(function() {
			delay(function() {
				$('#' + idTag + 'Toggle').trigger('click.bs.dropdown');
				quickSearch(url, idTag);
			}, 700);
		});
	/*
	$('#' + idTag + 'title')
		.mouseup(function() {
			delay(function() {
				$('#' + idTag + 'Toggle').trigger('click.bs.dropdown');
				quickSearch(url, idTag);
			}, 700);
		});
	*/
}


function setupDatePicker() {
	// $('.date-picker').datepicker({format: 'dd-mm-yyyy', startDate: '-2m',
	// endDate: '+1w', autoClose: true});
	$('.custom-datepicker').datepicker({format: "dd/mm/yyyy",autoclose: true});
}

function setupTinyMce(selector, mode) {
	tinymce.editors = [];
	tinymce.remove();

	if (mode) {
		tinymce.init({selector: '#' + selector, height: 300, readonly: mode});
	} else {
		tinymce.init({
			selector: '#' + selector,
			height: 300,
			plugins: ['code fullscreen fullpage']
			// plugins: 	[
			//   'advlist autolink lists link image charmap print preview
			//   anchor',
			//   'searchreplace visualblocks code fullscreen fullpage',
			//   'insertdatetime media table contextmenu paste code'
			// ],
			// toolbar: 'insertfile undo redo | styleselect | bold italic |
			// alignleft
			// aligncenter alignright alignjustify | bullist numlist outdent
			// indent |
			// link image'
		});
	}
}

function scrollTo(hash) {
	location.hash = "#" + hash;
}

function quickSearch(url, idTag) {
	if (idTag == "") {
		return;
	}

	var searchtext = $('#' + idTag + 'title').val();
	searchtext = searchtext.trim();

	var where = $('#' + idTag + 'where').val();
	

	$('#' + idTag + 'Dropdown')
		.html(
			'<li> &nbsp;&nbsp;&nbsp; <center><img src="../files/img/siteloader.gif"></center> &nbsp;&nbsp;&nbsp; </li>');
	$.post(url, "action=quicksearch&title=" + searchtext + "&tag=" + idTag +
					"&" + where,
		   function(response, status) {

			   if (typeof response == 'object') {
				   for (var layout in response) {
					   layoutID = layout.replace('.', '');
					   switch (layoutID) {
						   case 'error':
							   $.notify({message: response[layout]},
										{type: 'danger', timer: 5000});
							   break;

						   case 'alert':
							   $.notify({message: response[layout]},
										{type: 'info', timer: 5000});
							   break;

						   case 'quicksearch':
							   $('#' + idTag + 'Dropdown').html(response[layout]);
							   $('#' + idTag + 'title').focus();
							   break;
					   }
				   }
			   }
			  
		   });
}

function quickForm(formUrl) {
	aUrl = formUrl.split("?");
	$.post(aUrl[0],
		   aUrl[1], function(response, status) { handleResponse(response); });
}


var pageCurrent = 0;
var pageNavigator = [];
function getForm(formUrl) {
	if (formUrl == "logout") {
		window.location.href = "/logout";
	} else {
		startLoading();
		aUrl = formUrl.split("?");
		$.post(aUrl[0], aUrl[1], function(response, status) {
			handleResponse(response);
		});
	}
}

function handleResponse(response) {
	if (typeof response == 'object') {
		for (var layout in response) {
			layoutID = layout.replace('.', '');
			switch (layoutID) {
				case 'error':
					$.notify({message: response[layout]},
							 {type: 'danger', timer: 5000});
					break;

				case 'sticky':
					$.notify({message: response[layout]},
							 {type: 'danger', delay: 0});
					break;

				case 'startLoading':
						setTimeout(function() { startLoading() }, 600);
					break;

				case 'warning':
					warningMessage(response[layout]);
					break;

				case 'preview':
					// alert("We Previewing");
					// var site= "We Previewing!!";
					var site = response[layout];
					var popupWindow =
						window.open("", "", "menubar=0,scrollbars=0");
					popupWindow.document.write(site);
					popupWindow.document.close();
					break;

				case 'redirect':
					window.location.replace(response[layout]);
					break;


				case 'appSearchClearTag':
					clearActiveKeyword(response[layout]);
					break;

				case 'getform':
					getForm(response[layout]);
					break;

				case 'triggerSearch':
					$('.searchTrigger').trigger("submit");
					break;

				case 'triggerAppSearch':
					appSearchSubmit();
					break;

				case 'triggerSubSearch':
					$('.subsearchTrigger').trigger("submit");
					break;

				case 'toggleSubForm':
					toggleAppSidebar('subForm');
					toggleAppSidebar('subForm');
					break;

				case 'alert':
					$.notify({message: response[layout]},
							 {type: 'info', timer: 5000});
					break;

				case 'alertSuccess':
					$.notify({message: response[layout]},
							 {type: 'success', timer: 10000});
					break;

				case 'mainpanelContentSearch':
					$('#' + layoutID).html(response[layout]);
					closeView();
					break;

				case 'mainpanelContent':
					$('#' + layoutID).html(response[layout]);
					$('#mainpanelContent').show();
					$('#mainpanelContentSearch').hide();
					break;

				case 'searchresult': //appends search result for scrolling in mobile app version

					if ($('#appcontainer').attr("class") == "appcontainer") {
						if ($('#offset').val() == "0" ) {
							$('#' + layoutID).html(response[layout]);
						} else {					
							if(response[layout] == "") {
								$('#appendSearchDiv').css("min-height","0px");
								$('#appendSearchDiv').css("border","0");
								$('#appendSearchDiv').html("");
								$('#offset').val("0");
							} else {
								$('#' + layoutID).append(response[layout]);
								$('#appendSearchDiv').remove();
							}
						}
					} else {
						$('#' + layoutID).html(response[layout]);
					}

					break;

				default:
					$('#' + layoutID).html(response[layout]);
					// $('#'+layoutID).replaceWith(response[layout]);
					break;
			}
		}
	} else {
		$('html').html(response);
		// document.write(response); 
	}
	stopLoading();
}


function startLoading() {
	stopLoading();
	$("#pageContent").hide();
	$("#pageLoading").fadeIn(10);
}

function stopLoading() {
	$("#pageContent").fadeIn(10);
	$("#pageLoading").hide();
	adjustFooter();
}

function closeView(){
	$('#mainpanelContent').hide();
	$('#mainpanelContentSearch').show();
}

$('html')
	.on("click", ".btn-reset", function(event) { $('.resetForm').val('');$('.searchTrigger').trigger("submit"); });


function submitForm(formUrl, formData) {
	$.ajax({
		url: formUrl,
		type: 'POST',
		data: formData,
		async: true,
		cache: false,
		contentType: false,
		processData: false,
		success: function(response) { handleResponse(response); },
		error: function() { stopLoading(); }
	});
}

$('html').on("submit", ".form", function(event) {
	startLoading();
	event.preventDefault();
	var formData = new FormData($(this)[0]);
	var formUrl = $(this).attr('action') + "?";
	submitForm(formUrl, formData);
	return false;
});

$('html').on("focus", "#appSearchText", function(event) {
	$('#appSearchIcon').removeClass('pe-7s-close');$('#appSearchIcon').addClass('pe-7s-search');
});

$('html').on("click", "#appSearchIcon", function(event) {
	var curClass = $('#appSearchIcon').attr('class');
	if (curClass.indexOf('pe-7s-close') !== -1) {
		$('#appSearchText').val("");	
		$('#appSearchIcon').removeClass('pe-7s-close');
		$('#appSearchIcon').addClass('pe-7s-search');
	}


	if ($('#appSearchText').val() !== "" ) {
		$('#appSearchIcon').removeClass('pe-7s-search');
		$('#appSearchIcon').addClass('pe-7s-close');
	}

	$('.appSearchForm').trigger('submit');
});

$('html').on("submit", ".appSearchForm", function(event) {
	event.preventDefault();

	if ($('#appendSearchDiv').css("display") !== "table") {
		$('#offset').val("0");
	}

	if ($('#appSearchText').val() !== "" ) {
		$('#appSearchIcon').removeClass('pe-7s-search');
		$('#appSearchIcon').addClass('pe-7s-close');
	}
	
	var formData = new FormData($(this)[0]);
	var formUrl = $(this).attr('action') + "?";
	submitForm(formUrl, formData);
	return false;
});

// $('html').on("submit", "#form", function (event) {
// 	startLoading();
// 	event.preventDefault();
// 	var formData = new FormData($(this)[0]);
// 	var formUrl = $('#form').attr('action')+"?";
// 	submitForm(formUrl, formData)
// 	return false;
// });

$('html').on("submit", "#formDeactivateAll", function(event) {
	startLoading();
	event.preventDefault();
	var formData = new FormData($(this)[0]);
	var formUrl = $('#formDeactivateAll').attr('action') + "?";
	submitForm(formUrl, formData);
	return false;
});

$('html').on("submit", "#formSelected", function(event) {
	startLoading();
	event.preventDefault();
	var formData = new FormData($(this)[0]);
	var formUrl = $('#formSelected').attr('action') + "?";
	submitForm(formUrl, formData);
	return false;
});




// Disclaimer
function disclaimer(product) {
	$.notify(
		{
		  message:
			  "<b>Legal Disclaimer:</b> <br>Usage of <i>" + product +
				  "</i> without prior mutual consent is illegal. <br>It is the end user's responsibility to obey all applicable laws. <br>Developer assumes no liability and is not responsible for misuse."
		},
		{type: 'info', timer: 5000});
}

function error(message) {
	$.notify({message: message}, {type: 'danger', timer: 600});
}


function warningMessage(message) {
	$.notify({message: message}, {type: 'warning', timer: 400});
}

function successMessage(message) {
	$.notify({message: message}, {type: 'success', timer: 400});
}

function valuedMessage(message) {
	$.notify({message: message}, {type: 'warning', timer: 400});
}


// Title Alert
(function(a) {
	a.titleAlert = function(e, c) {
		if (a.titleAlert._running) {
			a.titleAlert.stop()
		}
		a.titleAlert._settings = c = a.extend({}, a.titleAlert.defaults, c);
		if (c.requireBlur && a.titleAlert.hasFocus) {
			return
		}
		c.originalTitleInterval = c.originalTitleInterval || c.interval;
		a.titleAlert._running = true;
		a.titleAlert._initialText = document.title;
		document.title = e;
		var b = true;
		var d = function() {
			if (!a.titleAlert._running) {
				return
			}
			b = !b;
			document.title = (b ? e : a.titleAlert._initialText);
			a.titleAlert._intervalToken =
				setTimeout(d, (b ? c.interval : c.originalTitleInterval))
		};
		a.titleAlert._intervalToken = setTimeout(d, c.interval);
		if (c.stopOnMouseMove) {
			a(document).mousemove(function(f) {
				a(this).unbind(f);
				a.titleAlert.stop()
			})
		}
		if (c.duration > 0) {
			a.titleAlert._timeoutToken =
				setTimeout(function() { a.titleAlert.stop() }, c.duration)
		}
	};
	a.titleAlert.defaults = {
		interval: 500,
		originalTitleInterval: null,
		duration: 0,
		stopOnFocus: true,
		requireBlur: false,
		stopOnMouseMove: false
	};
	a.titleAlert.stop = function() {
		clearTimeout(a.titleAlert._intervalToken);
		clearTimeout(a.titleAlert._timeoutToken);
		document.title = a.titleAlert._initialText;
		a.titleAlert._timeoutToken = null;
		a.titleAlert._intervalToken = null;
		a.titleAlert._initialText = null;
		a.titleAlert._running = false;
		a.titleAlert._settings = null
	};
	a.titleAlert.hasFocus = true;
	a.titleAlert._running = false;
	a.titleAlert._intervalToken = null;
	a.titleAlert._timeoutToken = null;
	a.titleAlert._initialText = null;
	a.titleAlert._settings = null;
	a.titleAlert._focus = function() {
		a.titleAlert.hasFocus = true;
		if (a.titleAlert._running && a.titleAlert._settings.stopOnFocus) {
			var b = a.titleAlert._initialText;
			a.titleAlert.stop();
			setTimeout(function() {
				if (a.titleAlert._running) {
					return
				}
				document.title = ".";
				document.title = b
			}, 1000)
		}
	};
	a.titleAlert._blur = function() { a.titleAlert.hasFocus = false };
	a(window).bind("focus", a.titleAlert._focus);
	a(window).bind("blur", a.titleAlert._blur)
})(jQuery);


// Maps Demo

demo = {
	initChartist: function() {

		var dataSales = {
			labels: ['Week 1', 'Week 2', 'Week 3', 'Week 4'],
			series: [
				[287, 385, 490, 492, 554, 586, 698, 695],
				[67, 152, 143, 240, 287, 335, 435, 437],
				[23, 113, 67, 108, 190, 239, 307, 308]
			]
		};

		var optionsSales = {
			lineSmooth: false,
			low: 0,
			high: 800,
			showArea: true,
			height: "245px",
			axisX: {
				showGrid: false,
			},
			lineSmooth: Chartist.Interpolation.simple({divisor: 3}),
			showLine: false,
			showPoint: false,
		};

		var responsiveSales = [[
			'screen and (max-width: 640px)',
			{
			  axisX: {
				  labelInterpolationFnc: function(value) { return value[0]; }
			  }
			}
		]];

		Chartist.Line('#chartHours', dataSales, optionsSales, responsiveSales);


		var data = {
			labels: [
				'Jan',
				'Feb',
				'Mar',
				'Apr',
				'Mai',
				'Jun',
				'Jul',
				'Aug',
				'Sep',
				'Oct',
				'Nov',
				'Dec'
			],
			series: [
				[542, 443, 320, 780, 553, 453, 326, 434, 568, 610, 756, 895],
				[412, 243, 280, 580, 453, 353, 300, 364, 368, 410, 636, 695],
				[421, 234, 203, 518, 457, 452, 321, 324, 318, 402, 616, 659]
			]
		};

		var options = {
			seriesBarDistance: 10,
			axisX: {showGrid: false},
			height: "245px"
		};

		var responsiveOptions = [[
			'screen and (max-width: 640px)',
			{
			  seriesBarDistance: 5,
			  axisX: {
				  labelInterpolationFnc: function(value) { return value[0]; }
			  }
			}
		]];

		Chartist.Bar('#chartActivity', data, options, responsiveOptions);

		var dataPreferences = {
			series: [[25, 30, 20, 25]]
		};

		var optionsPreferences = {
			donut: true,
			donutWidth: 40,
			startAngle: 0,
			total: 100,
			showLabel: false,
			axisX: {showGrid: false}
		};

		Chartist.Pie('#chartPreferences', dataPreferences, optionsPreferences);

		Chartist.Pie('#chartPreferences',
					 {labels: ['62%', '32%', '6%'], series: [62, 32, 6]});
	}
}
