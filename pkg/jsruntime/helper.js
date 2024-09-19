
function isUndefinedOrNull(value) {
    return value === undefined || value === null;
}

function hasUndefinedOrNull(value) {
    return Array.isArray(value)
        ? value.some(e => isUndefinedOrNull(e))
        : isUndefinedOrNull(value);
}

// Add custom Handlebars helpers below

Handlebars.registerHelper('replaceAll', function (target, old, newV) {
    if (hasUndefinedOrNull([target, old, newV])) {
        return target;
    }
    return target.replace(new RegExp(old, 'g'), newV);
});

Handlebars.registerHelper('match', function (pattern, target) {
    if (hasUndefinedOrNull([pattern, target])) {
        return false;
    }
    return new RegExp(pattern).test(target);
});

Handlebars.registerHelper('eq', function (left, right) {
    return left === right;
});

Handlebars.registerHelper('ne', function (left, right) {
    return left !== right;
});

Handlebars.registerHelper('snakeCase', function (s) {
    if (hasUndefinedOrNull(s)) {
        return s;
    }
    return s.replace(/([a-z])([A-Z])/g, '$1_$2').toLowerCase();
});

Handlebars.registerHelper('camelCase', function (s) {
    if (hasUndefinedOrNull(s)) {
        return s;
    }
    return s.replace(/([-_][a-z])/ig, function ($1) {
        return $1.toUpperCase()
            .replace('-', '')
            .replace('_', '');
    });
});

Handlebars.registerHelper('pascalCase', function (s) {
    if (hasUndefinedOrNull(s)) {
        return s;
    }
    return s.replace(/(\w)(\w*)/g, function ($0, $1, $2) {
        return $1.toUpperCase() + $2.toLowerCase();
    });
});

Handlebars.registerHelper('upperFirst', function (s) {
    if (hasUndefinedOrNull(s)) {
        return s;
    }
    return s.charAt(0).toUpperCase() + s.slice(1);
});

Handlebars.registerHelper('lowerFirst', function (s) {
    if (hasUndefinedOrNull(s)) {
        return s;
    }
    return s.charAt(0).toLowerCase() + s.slice(1);
});

Handlebars.registerHelper('uppercase', function (s) {
    if (hasUndefinedOrNull(s)) {
        return s;
    }
    return s.toUpperCase();
});

Handlebars.registerHelper('lowercase', function (s) {
    if (hasUndefinedOrNull(s)) {
        return s;
    }
    return s.toLowerCase();
});

Handlebars.registerHelper('trim', function (s) {
    if (hasUndefinedOrNull(s)) {
        return s;
    }
    return s.trim();
});

Handlebars.registerHelper('removePrefix', function (s, prefix) {
    if (hasUndefinedOrNull([s, prefix])) {
        return s;
    }
    return s.startsWith(prefix) ? s.slice(prefix.length) : s;
});

Handlebars.registerHelper('removeSuffix', function (s, suffix) {
    if (hasUndefinedOrNull([s, suffix])) {
        return s;
    }
    return s.endsWith(suffix) ? s.slice(0, -suffix.length) : s;
});
