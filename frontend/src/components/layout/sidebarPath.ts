export function isSidebarPathActive(currentPath: string, targetPath: string, knownPaths: readonly string[]): boolean {
  if (currentPath === targetPath) {
    return true
  }
  if (!currentPath.startsWith(targetPath + '/')) {
    return false
  }

  return !knownPaths.some((candidate) => (
    candidate !== targetPath &&
    candidate.length > targetPath.length &&
    candidate.startsWith(targetPath + '/') &&
    (currentPath === candidate || currentPath.startsWith(candidate + '/'))
  ))
}
