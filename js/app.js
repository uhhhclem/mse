var mseApp = angular.module('mseApp', ['ngMaterial']);

mseApp.controller('mseCtrl', function($scope, $http){
    $http.get('/api/newGame').success(function(d){
        $scope.board = d;
    });
})