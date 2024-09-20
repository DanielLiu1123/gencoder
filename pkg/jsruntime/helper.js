
function isUndefinedOrNull(value) {
    return value === undefined || value === null;
}

function hasUndefinedOrNull(value) {
    return Array.isArray(value)
        ? value.some(e => isUndefinedOrNull(e))
        : isUndefinedOrNull(value);
}

// Add custom Handlebars helpers below

Handlebars.registerHelper('_replaceAll', function (target, old, newV) {
    if (hasUndefinedOrNull([target, old, newV])) {
        return target;
    }
    return target.split(old).join(newV);
});

Handlebars.registerHelper('_match', function (pattern, target) {
    if (hasUndefinedOrNull([pattern, target])) {
        return false;
    }
    return new RegExp(pattern).test(target);
});

Handlebars.registerHelper('_eq', function (left, right) {
    return left === right;
});

Handlebars.registerHelper('_ne', function (left, right) {
    return left !== right;
});

Handlebars.registerHelper('_snakeCase', function (s) {
    if (hasUndefinedOrNull(s)) {
        return s;
    }
    return s.replace(/([a-z])([A-Z])/g, '$1_$2').toLowerCase();
});

Handlebars.registerHelper('_camelCase', function (s) {
    if (hasUndefinedOrNull(s)) {
        return s;
    }
    return s.replace(/([-_][a-z])/ig, function ($1) {
        return $1.toUpperCase()
            .replace('-', '')
            .replace('_', '');
    });
});

Handlebars.registerHelper('_pascalCase', function (s) {
    if (hasUndefinedOrNull(s)) {
        return s;
    }
    return s.replace(/(?:^|_)(\w)/g, function (_, $1) {
        return $1.toUpperCase();
    }).replace(/_/g, '');
});

Handlebars.registerHelper('_upperFirst', function (s) {
    if (hasUndefinedOrNull(s)) {
        return s;
    }
    return s.charAt(0).toUpperCase() + s.slice(1);
});

Handlebars.registerHelper('_lowerFirst', function (s) {
    if (hasUndefinedOrNull(s)) {
        return s;
    }
    return s.charAt(0).toLowerCase() + s.slice(1);
});

Handlebars.registerHelper('_uppercase', function (s) {
    if (hasUndefinedOrNull(s)) {
        return s;
    }
    return s.toUpperCase();
});

Handlebars.registerHelper('_lowercase', function (s) {
    if (hasUndefinedOrNull(s)) {
        return s;
    }
    return s.toLowerCase();
});

Handlebars.registerHelper('_trim', function (s) {
    if (hasUndefinedOrNull(s)) {
        return s;
    }
    return s.trim();
});

Handlebars.registerHelper('_removePrefix', function (s, prefix) {
    if (hasUndefinedOrNull([s, prefix])) {
        return s;
    }
    return s.startsWith(prefix) ? s.slice(prefix.length) : s;
});

Handlebars.registerHelper('_removeSuffix', function (s, suffix) {
    if (hasUndefinedOrNull([s, suffix])) {
        return s;
    }
    return s.endsWith(suffix) ? s.slice(0, -suffix.length) : s;
});
