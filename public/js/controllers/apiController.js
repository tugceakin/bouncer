/**
 * Created by tugceakin on 3/13/16.
 */


bouncerApp.controller('ApiController', function($scope, $http) {
    console.log('in controller');
    var url = "http://localhost:3000/profile";
    $scope.pageClass = 'page-page1';

    $scope.submit = function() {
        var data = $.param({
            json: JSON.stringify({
                item1: $scope.item1, item2: $scope.item2
            })
        });
        $http({
            method: 'POST',
            url: url,
            headers: {'Content-Type': 'application/x-www-form-urlencoded'}, //x-www-form-urlencoded
                        //'Access-Control-Allow-Origin': '*'},
            data: {item1: $scope.item1, item2: $scope.item2}
        }).success(function (data) {
            console.log(JSON.stringify(data));
        });
        console.log('clicked');

        //$http.post(url, data).success(function(data, status) {
        //    console.log(data);
        //})
    };
});