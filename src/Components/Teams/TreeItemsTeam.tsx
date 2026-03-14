import React, { useState, useEffect, useMemo, useCallback } from 'react';
import UserCardModal from './UserCardModal';
// ==================== Типы ====================

type NodeType = 'department' | 'user';

interface BaseNode {
  id: string;
  parentId: string | null;
  type: NodeType;
  name: string;
  email?: string,
  number?: string,
}

interface DepartmentNode extends BaseNode {
  type: 'department';
}

interface UserNode extends BaseNode {
  type: 'user';
  email?: string;
  number?: string;
  doljnost?: string;
}

type Node = DepartmentNode | UserNode | any;

interface TreeNode extends Node {
  children: TreeNode[];
}

// ==================== Хуки ====================

function useLocalStorage<T>(key: string, initialValue: T): [T, (value: T | ((val: T) => T)) => void] {
  const [storedValue, setStoredValue] = useState<T>(() => {
    try {
      return initialValue;      
//      const item = localStorage.getItem(key);
//      return item ? JSON.parse(item) : initialValue;
    } catch (error) {
      console.error(error);
      return initialValue;
    }
  });

  const setValue = (value: T | ((val: T) => T)) => {
    try {
      const valueToStore = value instanceof Function ? value(storedValue) : value;
      setStoredValue(valueToStore);
      localStorage.setItem(key, JSON.stringify(valueToStore));
    } catch (error) {
      console.error(error);
    }
  };

  return [storedValue, setValue];
}

function useDebounce<T>(value: T, delay: number): T {
  const [debouncedValue, setDebouncedValue] = useState<T>(value);

  useEffect(() => {
    const handler = setTimeout(() => {
      setDebouncedValue(value);
    }, delay);

    return () => {
      clearTimeout(handler);
    };
  }, [value, delay]);

  return debouncedValue;
}

// ==================== Функции работы с данными ====================

/**
 * Фильтрует плоский массив узлов, оставляя только те, которые соответствуют поисковому запросу,
 * а также всех их предков (чтобы сохранить структуру дерева).
 */
function filterTreeData(data: Node[], searchTerm: string): Node[] {
  if (!searchTerm.trim()) return data;

  const lowerSearch = searchTerm.toLowerCase();

  // Найти ID узлов, чьё имя содержит поисковую строку
  const matchingIds = new Set<string>();
  data.forEach(node => {
    if (node.name.toLowerCase().includes(lowerSearch) || 
        node.email?.toLowerCase().includes(lowerSearch) ||
        node.number?.toLowerCase().includes(lowerSearch) ) {
      matchingIds.add(node.id);
    }
  });

  // Добавить всех предков найденных узлов
  const idToNode = new Map<string, Node>(data.map(node => [node.id, node]));
  const allNeededIds = new Set<string>(matchingIds);

  matchingIds.forEach(id => {
    let currentId: string | null = id;
    while (currentId) {
      const node = idToNode.get(currentId);
      if (!node) break;
      allNeededIds.add(currentId);
      currentId = node.parentId;
    }
  });

  // Вернуть отфильтрованный массив
  return data.filter(node => allNeededIds.has(node.id));
}

/**
 * Строит дерево из плоского массива узлов.
 */
function buildTree(nodes: Node[], parentId: string | null = null): TreeNode[] {
  return nodes
    .filter(node => node.parentId === parentId)
    .map(node => ({
      ...node,
      children: buildTree(nodes, node.id),
    }));
}

// ==================== Компонент узла дерева ====================

interface TreeNodeProps {
  node: TreeNode;
  level?: number;
}

const TreeNodeComponent: React.FC<TreeNodeProps> = ({ node, level = 0 }) => {
  const [expanded, setExpanded] = useState(false);
  const hasChildren = node.children.length > 0;

  const toggleExpand = useCallback(() => setExpanded(prev => !prev), []);


  return (
    <>
    <li style={{ listStyle: 'none', marginLeft: level * 20 }}>
      <div
        style={{
          display: 'flex',
          alignItems: 'center',
          cursor: 'pointer',
          padding: '4px 0',
        }}
        onClick={toggleExpand}
      >
        {hasChildren && (
          <span style={{ width: 20, display: 'inline-block', textAlign: 'center' }}>
            {expanded ? '−' : '+'}
</span>
        )}
        {(node.type === 'department') ?
        <span style={{ marginLeft: hasChildren ? 0 : 20 }}>
          🏢 {node.name}
        </span>
        :
        <div className='row' style={{minWidth:"80%", marginLeft: hasChildren ? 0 : 20, textAlign:"left"}}>
        <div className='col'>
          👤 {node.name}
        </div>
        <div className='col'>
          {node.doljnost}
        </div>
        <div className='col' >
          {node.type === 'user' && node.email &&  
          <a href={`mailto:${node.email}?SUBJECT=Reactie`}>📧{node.email}</a>
          }
        </div>
        <div className='col'>
          {node.type === 'user' && node.number && node.number}
        </div>
        </div>
        }
      </div>
      {expanded && hasChildren && (
        <ul style={{ paddingLeft: 0, margin: 0 }}>
          {node.children.map(child => (
            <TreeNodeComponent key={child.id} node={child} level={level + 1} />
          ))}
        </ul>
      )}


    </li>
    </>
  );
};

// ==================== Основной компонент ====================

// Начальные данные (для примера)
const initialMockData: Node[] = [
  { id: '1', parentId: null, type: 'department', name: 'Большая компания' },
  { id: '2', parentId: '1', type: 'department', name: 'Управление' },
  { id: '3', parentId: '1', type: 'department', name: 'Службы 12' },
  { id: '4', parentId: '2', type: 'department', name: 'Служба 3' },
  { id: '5', parentId: '3', type: 'department', name: 'Поставки' },
  { id: '6', parentId: '2', type: 'department', name: 'Посредники' },
  { id: '7', parentId: '2', type: 'department', name: 'Перекупы' },
  { id: '8', parentId: '3', type: 'department', name: 'Разработка' },
  { id: '9', parentId: '1', type: 'department', name: 'Новости' },
  { id: '10', parentId: '1', type: 'department', name: 'Закупки' },
  { id: '11', parentId: '3', type: 'department', name: 'Самое длинное название какой-то канторы' },

  { id: '12', parentId: '2', type: 'user', name: 'Алексей', email: 'alex@example.com' },
  { id: '13', parentId: '2', type: 'user', name: 'Мария' },
  { id: '14', parentId: '5', type: 'user', name: 'Дмитрий' },
  { id: '15', parentId: '2', type: 'user', name: 'Алексей', email: 'alex@example.com', number:'📞12345'},
  { id: '16', parentId: '2', type: 'user', name: 'Мария' },
  { id: '17', parentId: '5', type: 'user', name: 'Дмитрий' },
  
  // Добавим ещё пользователей для демонстрации (около 100)
  ...Array.from({ length: 4000 }, (_, i) => ({
    id: `dev-${i}`,
    parentId: `${Math.floor(Math.random() * 11)}`,
    type: 'user' as const,
    name: `Роман ${i + 1}`,
    doljnost: `Оператор ${i + 1}`,
    email: `dev${i}@example.com`,
    number: `📞${Math.floor(Math.random() * 9)}${Math.floor(Math.random() * 9)}${Math.floor(Math.random() * 9)}${Math.floor(Math.random() * 9)}${Math.floor(Math.random() * 9)}`
  })),
];

const TreeItemsTeam: React.FC = () => {
  const [data, setData] = useLocalStorage<Node[]>('orgTree', []);
  const [loading, setLoading] = useState(false);
  const [searchTerm, setSearchTerm] = useState('');
  const debouncedSearch = useDebounce(searchTerm, 300);


  // Загрузка начальных данных, если localStorage пуст
  useEffect(() => {
    let isMounted = true;
    if (data.length === 0) {
      setLoading(true);
      // Имитация асинхронной загрузки с сервера
      setTimeout(() => {
        if (isMounted) {
          setData(initialMockData);
          setLoading(false);
        }
      }, 1000);
    }
    return () => {
      isMounted = false;
    };
  }, [data.length, setData]);

  const filteredData = useMemo(
    () => filterTreeData(data, debouncedSearch),
    [data, debouncedSearch]
  );



  const tree = useMemo(() => buildTree(filteredData), [filteredData]);

  return (
    <div style={{ padding: 20, fontFamily: 'sans-serif' }}>
      <input
        type="text"
        placeholder="Поиск по имени, email, номеру телефона..."
        value={searchTerm}
        onChange={e => setSearchTerm(e.target.value)}
        style={{
          marginBottom: 20,
          padding: 8,
          width: 600,
          borderRadius: 15,
          border: '1px solid #ccc',
        }}
      />
      {loading ? (
        <div style={{ textAlign: 'center', marginTop: 40, fontSize: 18 }}>Загрузка...</div>
      ) : (
        <ul style={{ paddingLeft: 0, margin: 0 }}>
          {tree.map(node => (
            <TreeNodeComponent key={node.id} node={node} />
          ))}
        </ul>
      )}

    </div>


  );
};

export default TreeItemsTeam;
